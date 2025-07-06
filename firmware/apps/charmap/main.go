package charmap

import (
	"eclair/hal/display"
	"eclair/hal/keypad"
	"eclair/hal/watchdog"
)

var pages = [][]byte{
	{'.', ',', '\'', '"', ':', ';', '-', '+'},
	{'=', '?', '!', '@', '#', '$', '%', '^'},
	{'&', '*', '(', ')', '[', ']', '<', '>'},
	{'{', '}', '\\', '/', '|', '_', '~', '`'},
	{0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87},
}

var disp *display.Display
var keys *keypad.Keypad

func refreshDisplay(page int) {
	disp.ClearBuffer()

	disp.DrawSprite16(icons, iconBack, 8, 0, display.AlignCenter, display.MaskNone, nil)
	disp.DrawSprite16(icons, iconNextPage, 8, 16, display.AlignCenter, display.MaskNone, nil)

	for chrId, chr := range pages[page] {
		x := 28*(chrId%4) + 36
		y := 16*(chrId/4) + 0

		// characters on the first page are too close to each other
		if page == 0 && chrId < 4 {
			y -= 1
		}
		if page == 0 && chrId >= 4 {
			y += 1
		}

		disp.DrawText([]byte{chr}, x, y, display.AlignCenter)
	}

	disp.Display()
}

func wrap(value int, min int, max int) int {
	if value < min {
		return max
	}
	if value > max {
		return min
	}
	return value
}

func Run() int {
	page := 0
	result := 0

	// - display -

	disp = display.Configure()

	// - keypad -

	handler := func(et keypad.EventType, id int) {
		if et.Released() {
			result = int(pages[page][id])
		}
	}

	keys = keypad.Configure()

	keys.SetHandlers([]func(keypad.EventType){
		func(et keypad.EventType) {
			if et.Released() {
				result = -1
			}
		},
		func(et keypad.EventType) { handler(et, 0) },
		func(et keypad.EventType) { handler(et, 1) },
		func(et keypad.EventType) { handler(et, 2) },
		func(et keypad.EventType) { handler(et, 3) },
		func(et keypad.EventType) {
			if et.Released() {
				page = wrap(page+1, 0, len(pages)-1)
				refreshDisplay(page)
			}
		},
		func(et keypad.EventType) { handler(et, 4) },
		func(et keypad.EventType) { handler(et, 5) },
		func(et keypad.EventType) { handler(et, 6) },
		func(et keypad.EventType) { handler(et, 7) },
		nil,
		nil,
		nil,
		nil,
		nil,
	})

	// - main loop -

	refreshDisplay(page)

	for result == 0 {
		watchdog.Feed()
		keys.Scan()
	}

	return result
}
