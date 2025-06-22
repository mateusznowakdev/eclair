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
	brownout.Configure()
	watchdog.Configure()
	clocks.Configure()

	if reset.IsSoftReset() {
		launcher.Run()
	} else {
		notes.Run(notes.DefaultName)
	}
}
