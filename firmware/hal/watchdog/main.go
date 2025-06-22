package watchdog

import "machine"

// Configure enables the watchdog peripheral, which resets the device if timer
// is not reset ("fed") for a defined amount of time.
func Configure() {
	_ = machine.Watchdog.Configure(machine.WatchdogConfig{TimeoutMillis: 4000})
	_ = machine.Watchdog.Start()
}

// Feed resets a watchdog timer, to indicate the app is healthy.
func Feed() {
	machine.Watchdog.Update()
}
