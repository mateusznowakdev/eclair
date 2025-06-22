package erase

import (
	"machine"
	"time"

	"eclair/hal/display"

	"github.com/mateusznowakdev/minifs"
)

func Run() {
	err := minifs.Format(machine.Flash)

	disp := display.Configure()
	disp.ClearBuffer()

	if err == nil {
		disp.DrawText([]byte("- done -"), display.Width/2, 8, display.AlignCenter)
	} else {
		disp.DrawText([]byte("- error -"), display.Width/2, 8, display.AlignCenter)
	}

	disp.Display()
	time.Sleep(2 * time.Second)

	return
}
