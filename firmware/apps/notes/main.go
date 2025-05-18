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

func deleteLine(note *Note) {
	lines := display.GetLines(note.file.Data)

	lineNo := getCursorLineNumber(lines, note.cursor)
	line := lines[lineNo]

	note.file.Data = slices.Delete(note.file.Data, line.Start, line.End)
	note.cursor = line.Start
	note.markDirty()
}

func dimScreen(disp *display.Display) {
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

func insertExtra(note *Note, disp *display.Display) {
	result := charmap.Run(disp)
	if result > 0 {
		note.insert(byte(result))
	}
}

func nextLine(note *Note) {
	lines := display.GetLines(note.file.Data)
	lineNo := getCursorLineNumber(lines, note.cursor)

	if lineNo < len(lines)-1 {
		// go down one line at the time
		note.cursor = lines[lineNo+1].Start
	} else if note.cursor < len(note.file.Data) {
		// then go to the end of a line
		note.cursor = len(note.file.Data)
	}
}

func prevLine(note *Note) {
	lines := display.GetLines(note.file.Data)

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

func refreshDisplay(disp display.Display, batt battery.Battery, note Note, shift bool) {
	disp.ClearBuffer()
	disp.DrawMultiText(note.file.Data, note.cursor)

	if note.dirty() {
		icon := icons[iconFile]
		disp.DrawSprite(icon, 0, uint(display.Width-len(icon)))
	} else if !batt.Good() {
		icon := icons[iconBattery]
		disp.DrawSprite(icon, 0, uint(display.Width-len(icon)))
	}

	if shift {
		icon := icons[iconShift]
		disp.DrawSprite(icon, 0, 0)
	}

	disp.Display()
}

func saveAndExit(note *Note) {
	_, err := note.write()
	if err != nil {
		reset.Lock()
	}

	reset.SoftReset()
}

func sendToPC(note *Note) {
	for _, char := range note.file.Data {
		keys := keymap[char-32]
		for _, key := range keys {
			keyboard.Keyboard.Down(key)
			time.Sleep(10 * time.Millisecond)
		}
		for _, key := range slices.Backward(keys) {
			keyboard.Keyboard.Up(key)
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

func Run() {
	shift := false

	note := NewNote()
	_, err := note.read()
	if err != nil {
		reset.Lock()
	}

	// - display -

	disp := display.New()

	// - keypad -

	keys := keypad.New()

	keys.SetBoolHandlers([]func(keypad.EventType) bool{
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { saveAndExit(&note) }, 'q', 'w', '1')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, nil, 'e', 'r', '2')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { prevLine(&note) }, 't', 'y', '3')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { insertExtra(&note, &disp) }, 'u', 'i', '4')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { note.delete() }, 'o', 'p', '5')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { shift = !shift }, 'a', 's', '6')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { note.cursorLeft() }, 'd', 'f', '7')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { nextLine(&note) }, 'g', 'h', '8')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { note.cursorRight() }, 'j', 'k', '9')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { sendToPC(&note) }, 'l', '-', '0')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, nil, 'z', 'x', '!')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { deleteLine(&note) }, 'c', 'v', '?')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { dimScreen(&disp) }, 'b', 'n', '\'')
		},
		func(et keypad.EventType) bool {
			return handler(&note, shift, et, func() { note.insert(' ') }, 'm', '.', ',')
		},
		nil,
	})

	// - battery -

	batt := battery.New()

	// - main loop -

	refreshDisplay(disp, batt, note, shift)

	for {
		watchdog.Feed()

		changed, err := note.writeDelayed()
		if err != nil {
			reset.Lock()
		}

		changed = changed || keys.Scan()
		if !changed {
			continue
		}

		refreshDisplay(disp, batt, note, shift)
	}
}
