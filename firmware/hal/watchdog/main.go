package watchdog

import (
	"device/sam"

	"eclair/hal/clocks"
)

// Configure enables the watchdog peripheral, which resets the device if timer
// is not reset ("fed") for a defined amount of time.
func Configure() {
	clocks.ConfigureWatchdog()

	// power on the WDT peripheral
	sam.PM.APBAMASK.SetBits(sam.PM_APBAMASK_WDT_)

	// disable peripheral
	sam.WDT.CTRL.ClearBits(sam.WDT_CTRL_ENABLE)
	for sam.WDT.STATUS.HasBits(sam.WDT_STATUS_SYNCBUSY) {
	}

	// set 4096 cycles (approximately 4 seconds)
	sam.WDT.CONFIG.Set(sam.WDT_CONFIG_PER_4K)

	// enable peripheral
	sam.WDT.CTRL.SetBits(sam.WDT_CTRL_ENABLE)
	for sam.WDT.STATUS.HasBits(sam.WDT_STATUS_SYNCBUSY) {
	}
}

// Feed resets a watchdog timer, to indicate the app is healthy.
func Feed() {
	sam.WDT.CLEAR.Set(sam.WDT_CLEAR_CLEAR_KEY) // other values reset the CPU immediately
}
