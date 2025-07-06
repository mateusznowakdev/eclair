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

var disp *display.Display
var keys *keypad.Keypad

const acceleration = 1.1
const maxSpeed = 8.0

func Run() {
	speedX := float64(0)
	speedY := float64(0)

	// - display -

	disp = display.Configure()

	disp.ClearBuffer()
	disp.DrawSprite16(icons, iconBack, 8, 0, display.AlignCenter, display.MaskNone, nil)
	disp.DrawText([]byte("T,D,G,J: Move"), display.Width/2+8, 0, display.AlignCenter)
	disp.DrawText([]byte("C,B: Click"), display.Width/2+8, 16, display.AlignCenter)
	disp.Display()

	disp.SetContrast(display.ContrastLow)

	// - keypad -

	keys = keypad.Configure()

	keys.SetHandlers([]func(keypad.EventType){
		func(et keypad.EventType) {
			if et.Released() {
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
		watchdog.Feed()
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
