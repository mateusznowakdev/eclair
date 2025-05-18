package charmap

import (
	"eclair/hal/display"
	"eclair/hal/keypad"
	"eclair/hal/watchdog"
)

var chars = [][]byte{
	{'.', ',', ';', ':', '\'', '"', '-', '+'},
	{'=', '?', '!', '@', '#', '$', '%', '^'},
	{'&', '*', '(', ')', '[', ']', '<', '>'},
	{'{', '}', '\\', '/', '|', '_', '~', '`'},
}

// there is no support for horizontal text centering, and having one would
// slow down the main text renderer, so for now have this offset map instead
var offs = [][]uint8{
	{5, 5, 5, 5, 5, 3, 2, 2},
	{2, 3, 5, 1, 1, 2, 2, 3},
	{2, 2, 3, 3, 4, 4, 3, 3},
	{2, 3, 2, 3, 5, 2, 2, 3},
}

func refreshDisplay(disp *display.Display, posX int, posY int) {
	disp.ClearBufferTop()

	for charNo, char := range chars[posY] {
		disp.DrawText([]byte{char}, 0, uint(charNo)*16+uint(offs[posY][charNo]))
	}

	disp.DrawTextFrame(0, uint(posX*16), uint(posX*16+14))
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

func Run(disp *display.Display) int {
	posX := 0
	posY := 0

	result := 0

	// - using existing display instance to prevent it from being reset -

	// - keypad -

	keys := keypad.New()

	keys.SetHandlers([]func(keypad.EventType){
		nil,
		nil,
		func(et keypad.EventType) {
			if et.Released() {
				posY = wrap(posY-1, 0, len(chars)-1)
				refreshDisplay(disp, posX, posY)
			}
		},
		nil,
		nil,
		nil,
		func(et keypad.EventType) {
			if et.Released() {
				posX = wrap(posX-1, 0, len(chars[0])-1)
				refreshDisplay(disp, posX, posY)
			}
		},
		func(et keypad.EventType) {
			if et.Released() {
				posY = wrap(posY+1, 0, len(chars)-1)
				refreshDisplay(disp, posX, posY)
			}
		},
		func(et keypad.EventType) {
			if et.Released() {
				posX = wrap(posX+1, 0, len(chars[0])-1)
				refreshDisplay(disp, posX, posY)
			}
		},
		func(et keypad.EventType) {
			if et.Released() {
				result = int(chars[posY][posX])
			}
		},
		nil,
		nil,
		nil,
		nil,
		nil,
	})

	// - main loop -

	refreshDisplay(disp, posX, posY)

	for result == 0 {
		watchdog.Feed()
		keys.Scan()
	}

	return result
}
