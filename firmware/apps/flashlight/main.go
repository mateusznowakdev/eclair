package flashlight

import (
	"eclair/display"
	"eclair/keypad"
	"eclair/peripherals"
)

func Run() {
	// - display -

	disp := display.NewDisplay()

	disp.ClearBuffer()
	disp.DrawText([]byte("- hold any key -"), 1, 10)
	disp.Display()

	disp.SetContrast(display.ContrastLow)

	// - keypad -

	handler := func(et keypad.EventType) {
		if et.Pressed() {
			disp.SetContrast(display.ContrastHigh)
			disp.SetInverted(true)
		} else {
			disp.SetContrast(display.ContrastLow)
			disp.SetInverted(false)
		}
	}

	keys := keypad.NewKeypad()

	keys.SetHandlers([]func(keypad.EventType){
		func(et keypad.EventType) {
			if et.Alt() && et.Released() {
				peripherals.SoftReset()
			}
			handler(et)
		},
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		nil,
	})

	// - main loop -

	for {
		peripherals.FeedWatchdog()
		keys.Scan()
	}
}
