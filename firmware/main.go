package main

import (
	"eclair/apps/launcher"
	"eclair/apps/notes"
	"eclair/peripherals"
)

func main() {
	peripherals.ConfigureBOD33()
	peripherals.ConfigureWatchdog()
	peripherals.ConfigureCPUClock()
	peripherals.ConfigureUSBClock()

	if peripherals.IsSoftReset() {
		launcher.Run()
	} else {
		notes.Run()
	}
}
