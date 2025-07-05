package launcher

import (
	"time"

	"eclair/hal/display"
	"eclair/hal/keypad"
	"eclair/hal/reset"
	"eclair/hal/watchdog"
)

var disp *display.Display
var keys *keypad.Keypad

type Entry struct {
	name       string
	entrypoint func()
}

func clamp(value int, min int, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func refreshDisplay(pos int) {
	disp.ClearBuffer()

	opt := apps[pos]
	disp.DrawText([]byte(opt.name), 2, 0, display.AlignLeft)
	//disp.DrawTextFrame(0, 126, 0)

	if pos < len(apps)-1 {
		opt = apps[pos+1]
		disp.DrawText([]byte(opt.name), 2, 16, display.AlignLeft)
	}

	disp.Display()
}

func Run() {
	pos := 0

	// - display -

	disp = display.Configure()

	// - keypad -

	keys = keypad.Configure()

	keys.SetHandlers([]func(keypad.EventType){
		nil,
		nil,
		func(et keypad.EventType) {
			if et.Released() {
				pos = clamp(pos-1, 0, len(apps)-1)
				refreshDisplay(pos)
			}
		},
		nil,
		nil,
		nil,
		nil,
		func(et keypad.EventType) {
			if et.Released() {
				pos = clamp(pos+1, 0, len(apps)-1)
				refreshDisplay(pos)
			}
		},
		nil,
		func(et keypad.EventType) {
			if et.Released() {
				disp.ClearDisplay()
				time.Sleep(250 * time.Millisecond)

				apps[pos].entrypoint()
				reset.SoftReset()
			}
		},
		nil,
		nil,
		nil,
		nil,
		nil,
	})

	// - main loop -

	refreshDisplay(pos)

	for {
		watchdog.Feed()
		keys.Scan()
	}
}
