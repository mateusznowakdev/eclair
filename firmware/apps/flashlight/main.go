package flashlight

import (
	"eclair/hal/display"
	"eclair/hal/keypad"
	"eclair/hal/reset"
	"eclair/hal/watchdog"
)

func Run() {
	// - display -

	disp := display.New()

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

	keys := keypad.New()

	keys.SetHandlers([]func(keypad.EventType){
		func(et keypad.EventType) {
			if et.Alt() && et.Released() {
				reset.SoftReset()
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
		watchdog.Feed()
		keys.Scan()
	}
}
