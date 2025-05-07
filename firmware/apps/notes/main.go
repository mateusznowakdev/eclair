package notes

import (
	"machine/usb/hid/keyboard"
	"slices"

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

func handler(et keypad.EventType, note *Note, alt func(), opts ...byte) bool {
	if alt != nil && et.Alt() && et.Released() {
		alt()
		return true
	}

	if et.Double() && et.Released() {
		for optNo, opt := range opts {
			if note.last() != opt {
				continue
			}

			optNo += 1
			if optNo == len(opts) {
				optNo = 0
			}

			note.replace(opts[optNo])
			return true
		}

		note.insert(opts[0])
		return true
	}

	if et.Released() {
		note.insert(opts[0])
		return true
	}

	return false
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

func refreshDisplay(disp display.Display, batt battery.Battery, note Note) {
	disp.ClearBuffer()

	if len(note.file.Data) > 0 {
		disp.DrawMultiText(note.file.Data, note.cursor)
	} else {
		disp.DrawText([]byte("- start typing -"), 1, 10)
	}

	if note.dirty() {
		disp.DrawSprite(icons["file"], 0, 0)
	} else if !batt.Good() {
		disp.DrawSprite(icons["battery"], 0, 0)
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
	_, _ = keyboard.Keyboard.Write(note.file.Data)
}

func Run() {
	note := NewNote()

	_, err := note.read()
	if err != nil {
		reset.Lock()
	}

	// - display -

	disp := display.NewDisplay()

	// - keypad -

	keys := keypad.NewKeypad()

	keys.SetBoolHandlers([]func(keypad.EventType) bool{
		func(et keypad.EventType) bool {
			return handler(et, &note, func() { saveAndExit(&note) }, 'q', 'w', '1')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, nil, 'e', 'r', '2')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, func() { prevLine(&note) }, 't', 'y', '3')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, nil, 'u', 'i', '4')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, func() { note.delete() }, 'o', 'p', '5')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, nil, 'a', 's', '6')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, func() { note.cursorLeft() }, 'd', 'f', '7')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, func() { nextLine(&note) }, 'g', 'h', '8')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, func() { note.cursorRight() }, 'j', 'k', '9')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, func() { sendToPC(&note) }, 'l', '-', '0')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, nil, 'z', 'x', '!')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, func() { deleteLine(&note) }, 'c', 'v', '?')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, func() { dimScreen(&disp) }, 'b', 'n', '\'')
		},
		func(et keypad.EventType) bool {
			return handler(et, &note, func() { note.insert(' ') }, 'm', '.', ',')
		},
		nil,
	})

	// - battery -

	batt := battery.NewBattery()

	// - main loop -

	refreshDisplay(disp, batt, note)

	for {
		watchdog.FeedWatchdog()

		changed, err := note.writeDelayed()
		if err != nil {
			reset.Lock()
		}

		changed = changed || keys.Scan()
		if !changed {
			continue
		}

		refreshDisplay(disp, batt, note)
	}
}
