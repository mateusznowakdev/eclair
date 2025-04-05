package launcher

import (
	"time"

	"eclair/display"
	"eclair/keypad"
	"eclair/peripherals"
)

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

func refreshDisplay(disp display.Display, pos int) {
	disp.ClearBuffer()

	opt := apps[pos]
	disp.DrawText([]byte(">"), 0, 0)
	disp.DrawText([]byte(opt.name), 0, 10)

	if pos < len(apps)-1 {
		opt = apps[pos+1]
		disp.DrawText([]byte(opt.name), 2, 10)
	}

	disp.Display()
}

func Run() {
	pos := 0

	// - display -

	disp := display.NewDisplay()

	// - keypad -

	keys := keypad.NewKeypad()

	keys.SetHandlers([]func(keypad.EventType){
		nil,
		nil,
		func(et keypad.EventType) {
			if et.Released() {
				pos = clamp(pos-1, 0, len(apps)-1)
				refreshDisplay(disp, pos)
			}
		},
		nil,
		nil,
		nil,
		nil,
		func(et keypad.EventType) {
			if et.Released() {
				pos = clamp(pos+1, 0, len(apps)-1)
				refreshDisplay(disp, pos)
			}
		},
		nil,
		func(et keypad.EventType) {
			if et.Released() {
				disp.ClearDisplay()
				time.Sleep(250 * time.Millisecond)

				apps[pos].entrypoint()
			}
		},
		nil,
		nil,
		nil,
		nil,
		nil,
	})

	// - main loop -

	refreshDisplay(disp, pos)

	for {
		peripherals.FeedWatchdog()
		keys.Scan()
	}
}
