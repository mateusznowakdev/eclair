package erase

import (
	"machine"
	"time"

	"eclair/hal/display"
	"eclair/hal/keypad"
	"eclair/hal/reset"
	"eclair/hal/watchdog"

	"github.com/mateusznowakdev/minifs"
)

var disp *display.Display
var keys *keypad.Keypad

func refreshDisplay(state int, seconds int64) {
	disp.ClearBuffer()
	disp.DrawSprite16(icons, iconBack, 8, 0, display.AlignCenter, display.MaskNone, nil)

	x := display.Width/2 + 8

	switch state {
	case 2:
		disp.DrawSprite16(icons, iconFailure, x, 0, display.AlignCenter, display.MaskNone, nil)
		disp.DrawText([]byte("Check console"), x, 16, display.AlignCenter)
	case 1:
		disp.DrawSprite16(icons, iconSuccess, x, 0, display.AlignCenter, display.MaskNone, nil)
		disp.DrawText([]byte("Done"), x, 16, display.AlignCenter)
	default:
		text := []byte{byte('0' + seconds)}
		disp.DrawText(text, x, 8, display.AlignCenter)
	}

	disp.Display()
}

func Run() {
	// - display -

	disp = display.Configure()

	// - keypad -

	keys = keypad.Configure()

	keys.SetHandlers([]func(keypad.EventType){
		func(et keypad.EventType) {
			if et.Released() {
				reset.SoftReset()
			}
		},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	})

	// - main loop -

	timeout := time.Now().Add(3 * time.Second).UnixMilli()
	state := 0

	for {
		now := time.Now().UnixMilli()

		if state == 0 {
			if now >= timeout {
				err := minifs.Format(machine.Flash)
				if err == nil {
					state = 1
				} else {
					state = 2
					print(err, "\r\n")
				}
			}

			refreshDisplay(state, (timeout-now)/1000+1)
			time.Sleep(20 * time.Millisecond)
		}

		watchdog.Feed()
		keys.Scan()
	}
}
