package clocks

import "device/sam"

const clockGCLK0 = 8_000_000 // Hz
const clockGCLK0Divider = 48_000_000 / clockGCLK0

// ConfigureCPU changes the GCLK0 clock frequency to 8MHz. This is the main
// clock and is used for CPU and most peripherals.
//
// This will make the USB peripheral unusable, so ConfigureUSB should be called
// immediately.
func ConfigureCPU() {
	sam.GCLK.GENDIV.Set(0 | (clockGCLK0Divider << sam.GCLK_GENDIV_DIV_Pos))
	for sam.GCLK.STATUS.HasBits(sam.GCLK_STATUS_SYNCBUSY) {
	}
}

// ConfigureUSB enables the GCLK5 clock at 48MHz and assigns it to the USB
// peripheral.
func ConfigureUSB() {
	sam.GCLK.GENDIV.Set(5 | (1 << sam.GCLK_GENDIV_DIV_Pos))
	sam.GCLK.GENCTRL.Set(5 | (sam.GCLK_GENCTRL_SRC_DFLL48M << sam.GCLK_GENCTRL_SRC_Pos) | sam.GCLK_GENCTRL_GENEN)
	for sam.GCLK.STATUS.HasBits(sam.GCLK_STATUS_SYNCBUSY) {
	}

	sam.GCLK.CLKCTRL.Set(sam.GCLK_CLKCTRL_ID_USB | (sam.GCLK_CLKCTRL_GEN_GCLK5 << sam.GCLK_CLKCTRL_GEN_Pos) | sam.GCLK_CLKCTRL_CLKEN)
}

// ConfigureWatchdog enables the GCLK8 clock at 1MHz and assigns it to the
// watchdog peripheral.
func ConfigureWatchdog() {
	sam.GCLK.GENDIV.Set(8 | (32 << sam.GCLK_GENDIV_DIV_Pos))
	sam.GCLK.GENCTRL.Set(8 | (sam.GCLK_GENCTRL_SRC_OSCULP32K << sam.GCLK_GENCTRL_SRC_Pos) | sam.GCLK_GENCTRL_GENEN)
	for sam.GCLK.STATUS.HasBits(sam.GCLK_STATUS_SYNCBUSY) {
	}

	sam.GCLK.CLKCTRL.Set(sam.GCLK_CLKCTRL_ID_WDT | (sam.GCLK_CLKCTRL_GEN_GCLK8 << sam.GCLK_CLKCTRL_GEN_Pos) | sam.GCLK_CLKCTRL_CLKEN)
}

// PatchedGCLK0Frequency returns a valid frequency for the SPI peripheral, based
// on the custom GCLK0 prescaler value. This workaround is needed because
// machine.CPUFrequency on SAMD21 returns a hardcoded value of 48MHz.
func PatchedGCLK0Frequency(value int) int {
	return value * clockGCLK0Divider
}
