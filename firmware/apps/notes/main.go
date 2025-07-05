package notes

import (
	"machine/usb/hid/keyboard"
	"slices"
	"time"

	"eclair/apps/charmap"
	"eclair/hal/battery"
	"eclair/hal/display"
	"eclair/hal/keypad"
	"eclair/hal/reset"
	"eclair/hal/watchdog"
)

var batt *battery.Battery
var disp *display.Display
var keys *keypad.Keypad

func deleteLine(note *Note) {
	lines := display.GetLines(note.data)

	lineNo := getCursorLineNumber(lines, note.cursor)
	line := lines[lineNo]

	note.data = slices.Delete(note.data, line.Start, line.End)
	note.cursor = line.Start
	note.markDirty()
}

func dimScreen() {
	if disp.Contrast() != display.ContrastLow {
		disp.SetContrast(display.ContrastLow)
	} else {
		disp.SetContrast(display.ContrastNormal)
	}
}

func getCursorLineNumber(lines []display.Line, cursor int) int {
	for ln := len(lines) - 1; ln >= 0; ln-- {
		line := lines[ln]
		if cursor >= line.Start {
			return ln
		}
	}

	return 0
}

func handler(note *Note, shift bool, et keypad.EventType, alt func(), opts ...byte) bool {
	if alt != nil && et.Alt() && et.Released() {
		alt()
		return true
	}

	if et.Double() && et.Released() {
		last := upper(note.last(), shift)

		for optNo, opt := range opts {
			opt = upper(opt, shift)
			if last != opt {
				continue
			}

			optNo += 1
			if optNo == len(opts) {
				optNo = 0
			}

			note.replace(upper(opts[optNo], shift))
			return true
		}

		note.insert(upper(opts[0], shift))
		return true
	}

	if et.Released() {
		note.insert(upper(opts[0], shift))
		return true
	}

	return false
}

func insertExtra(note *Note) {
	result := charmap.Run(disp)
	if result > 0 {
		note.insert(byte(result))
	}
}

func nextLine(note *Note) {
	lines := display.GetLines(note.data)
	lineNo := getCursorLineNumber(lines, note.cursor)

	if lineNo < len(lines)-1 {
		// go down one line at the time
		note.cursor = lines[lineNo+1].Start
	} else if note.cursor < len(note.data) {
		// then go to the end of a line
		note.cursor = len(note.data)
	}
}

func prevLine(note *Note) {
	lines := display.GetLines(note.data)

	lineNo := getCursorLineNumber(lines, note.cursor)
	line := lines[lineNo]

	if note.cursor > line.Start {
		// go to the beginning of a line
		note.cursor = line.Start
	} else if lineNo > 0 {
		// then go up one line at the time
		note.cursor = lines[lineNo-1].Start
	}
}

func refreshDisplay(note *Note, shift bool) {
	disp.ClearBuffer()
	disp.DrawMultiText(note.data, note.cursor)

	if note.dirty() {
		disp.DrawSprite8(icons, iconFile, display.Width-8, 0, display.MaskAll, nil)
	} else if !batt.Good() {
		disp.DrawSprite8(icons, iconBattery, display.Width-8, 0, display.MaskAll, nil)
	}

	if shift {
		disp.DrawSprite8(icons, iconShift, 0, 0, display.MaskAll, nil)
	}

	disp.Display()
}

func saveAndExit(note *Note) {
	err := note.write()
	if err != nil {
		panic(err)
	}

	reset.SoftReset()
}

func sendToPC(note *Note) {
	for _, char := range note.data {
		keys := keymap[char-32]
		for _, key := range keys {
			_ = keyboard.Keyboard.Down(key)
			time.Sleep(10 * time.Millisecond)
		}
		for _, key := range slices.Backward(keys) {
			_ = keyboard.Keyboard.Up(key)
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func upper(value byte, shift bool) byte {
	if shift && value >= 'a' && value <= 'z' {
		return value - 32
	}
	return value
}

func Run(name string) {
	var note *Note
	shift := false

	// - display -

	disp = display.Configure()

	// - keypad -

	keys = keypad.Configure()

	keys.SetBoolHandlers([]func(keypad.EventType) bool{
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { saveAndExit(note) }, 'q', 'w', '1')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, nil, 'e', 'r', '2')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { prevLine(note) }, 't', 'y', '3')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { insertExtra(note) }, 'u', 'i', '4')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { note.delete() }, 'o', 'p', '5')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { shift = !shift }, 'a', 's', '6')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { note.cursorLeft() }, 'd', 'f', '7')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { nextLine(note) }, 'g', 'h', '8')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { note.cursorRight() }, 'j', 'k', '9')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { sendToPC(note) }, 'l', '-', '0')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, nil, 'z', 'x', '!')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { deleteLine(note) }, 'c', 'v', '?')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { dimScreen() }, 'b', 'n', '\'')
		},
		func(et keypad.EventType) bool {
			return handler(note, shift, et, func() { note.insert(' ') }, 'm', '.', ',')
		},
		nil,
	})

	// - battery -

	batt = battery.Configure()

	// - load data -

	note, err := NewNote(name)
	if err != nil {
		disp.ClearBuffer()
		disp.DrawText([]byte("- no format -"), display.Width/2, 8, display.AlignCenter)
		disp.Display()

		time.Sleep(2 * time.Second)
		return
	}

	err = note.read()
	if err != nil {
		panic(err)
	}

	// - main loop -

	refreshDisplay(note, shift)

	for {
		watchdog.Feed()

		changed, err := note.writeDelayed()
		if err != nil {
			panic(err)
		}

		changed = changed || keys.Scan()
		if !changed {
			continue
		}

		refreshDisplay(note, shift)
	}
}
