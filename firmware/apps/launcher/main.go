package launcher

import (
	"time"

	"eclair/apps/bootloader"
	"eclair/hal/display"
	"eclair/hal/keypad"
	"eclair/hal/reset"
	"eclair/hal/watchdog"
)

var disp *display.Display
var keys *keypad.Keypad

type Entry struct {
	icon       int
	entrypoint func()
}

func Run() {
	// - display -

	disp = display.Configure()
	disp.ClearDisplay()

	for appId, app := range apps {
		if app == nil {
			continue
		}

		x := 28*(appId%5) + 2
		y := 20*(appId/5) - 2
		disp.DrawSprite16(icons, app.icon, x, y, display.MaskNone, nil)
	}

	disp.Display()

	// - keypad -

	handler := func(et keypad.EventType, id int) {
		if et.Released() && apps[id] != nil {
			disp.ClearDisplay()
			time.Sleep(250 * time.Millisecond)

			apps[id].entrypoint()
			reset.SoftReset()
		}
	}

	keys = keypad.Configure()

	keys.SetHandlers([]func(keypad.EventType){
		func(et keypad.EventType) { handler(et, 0) },
		func(et keypad.EventType) { handler(et, 1) },
		func(et keypad.EventType) { handler(et, 2) },
		func(et keypad.EventType) { handler(et, 3) },
		func(et keypad.EventType) { handler(et, 4) },
		func(et keypad.EventType) { handler(et, 5) },
		func(et keypad.EventType) { handler(et, 6) },
		func(et keypad.EventType) { handler(et, 7) },
		func(et keypad.EventType) { handler(et, 8) },
		func(et keypad.EventType) { handler(et, 9) },
		nil,
		nil,
		func(et keypad.EventType) {
			if et.Alt() && et.Released() {
				bootloader.Run()
			}
		},
		nil,
	})

	// - main loop -

	for {
		watchdog.Feed()
		keys.Scan()
	}
}
