package main

import (
	"eclair/apps/launcher"
	"eclair/apps/notes"
	"eclair/hal/brownout"
	"eclair/hal/clocks"
	"eclair/hal/reset"
	"eclair/hal/watchdog"
)

func main() {
	brownout.ConfigureBOD33()
	watchdog.ConfigureWatchdog()
	clocks.ConfigureCPUClock()
	clocks.ConfigureUSBClock()

	if reset.IsSoftReset() {
		launcher.Run()
	} else {
		notes.Run()
	}
}
