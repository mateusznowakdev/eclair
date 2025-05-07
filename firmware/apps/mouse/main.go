package mouse

import (
	"machine/usb/hid/mouse"
	"math"
	"time"

	"eclair/hal/display"
	"eclair/hal/keypad"
	"eclair/hal/reset"
	"eclair/hal/watchdog"
)

const acceleration = 1.1
const maxSpeed = 8.0

func Run() {
	speedX := float64(0)
	speedY := float64(0)

	// - display -

	disp := display.NewDisplay()

	disp.ClearBuffer()
	disp.DrawText([]byte("- mouse mode -"), 1, 10)
	disp.Display()

	disp.SetContrast(display.ContrastLow)

	// - keypad -

	keys := keypad.NewKeypad()

	keys.SetHandlers([]func(keypad.EventType){
		func(et keypad.EventType) {
			if et.Alt() && et.Released() {
				reset.SoftReset()
			}
		},
		nil,
		func(et keypad.EventType) {
			if et.Pressed() {
				speedY = -1
			} else {
				speedY = 0
			}
		},
		nil,
		nil,
		nil,
		func(et keypad.EventType) {
			if et.Pressed() {
				speedX = -1
			} else {
				speedX = 0
			}
		},
		func(et keypad.EventType) {
			if et.Pressed() {
				speedY = 1
			} else {
				speedY = 0
			}
		},
		func(et keypad.EventType) {
			if et.Pressed() {
				speedX = 1
			} else {
				speedX = 0
			}
		},
		nil,
		nil,
		func(et keypad.EventType) {
			if et.Pressed() {
				mouse.Mouse.Press(mouse.Left)
			} else {
				mouse.Mouse.Release(mouse.Left)
			}
		},
		func(et keypad.EventType) {
			if et.Pressed() {
				mouse.Mouse.Press(mouse.Right)
			} else {
				mouse.Mouse.Release(mouse.Right)
			}
		},
		nil,
		nil,
	})

	// - main loop -

	for {
		watchdog.FeedWatchdog()
		keys.Scan()

		mouse.Mouse.Move(int(speedX), int(speedY))

		if speedX != 0 && math.Abs(speedX) < maxSpeed {
			speedX = speedX * acceleration
		}
		if speedY != 0 && math.Abs(speedY) < maxSpeed {
			speedY = speedY * acceleration
		}

		time.Sleep(15 * time.Millisecond)
	}
}
