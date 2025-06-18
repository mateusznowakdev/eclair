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
	{0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87},
}

var disp *display.Display
var keys *keypad.Keypad

func refreshDisplay(posX int, posY int) {
	disp.ClearBufferTop()

	for charNo, char := range chars[posY] {
		disp.DrawText([]byte{char}, charNo*16+8, 0, display.AlignCenter)
	}

	disp.DrawTextFrame(uint(posX*16), uint(posX*16+14), 0)
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

func Run(dispParent *display.Display) int {
	posX := 0
	posY := 0

	result := 0

	// - using existing display instance to prevent it from being reset -

	disp = dispParent

	// - keypad -

	keys = keypad.Configure()

	keys.SetHandlers([]func(keypad.EventType){
		func(et keypad.EventType) {
			if et.Alt() && et.Released() {
				result = -1
			}
		},
		nil,
		func(et keypad.EventType) {
			if et.Released() {
				posY = wrap(posY-1, 0, len(chars)-1)
				refreshDisplay(posX, posY)
			}
		},
		nil,
		nil,
		nil,
		func(et keypad.EventType) {
			if et.Released() {
				posX = wrap(posX-1, 0, len(chars[0])-1)
				refreshDisplay(posX, posY)
			}
		},
		func(et keypad.EventType) {
			if et.Released() {
				posY = wrap(posY+1, 0, len(chars)-1)
				refreshDisplay(posX, posY)
			}
		},
		func(et keypad.EventType) {
			if et.Released() {
				posX = wrap(posX+1, 0, len(chars[0])-1)
				refreshDisplay(posX, posY)
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

	refreshDisplay(posX, posY)

	for result == 0 {
		watchdog.Feed()
		keys.Scan()
	}

	return result
}
