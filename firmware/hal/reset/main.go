package reset

import (
	"device/sam"
	"machine"
	"time"
)

// IsSoftReset checks if the CPU reset was caused by user action such as
// manually triggering the reset or leaving the bootloader mode.
func IsSoftReset() bool {
	return sam.PM.GetRCAUSE_SYST() != 0
}

// Lock waits for a watchdog to reset the device.
func Lock() {
	//goland:noinspection ALL
	for {
	}
}

// SoftReset performs a CPU reset, after a slight delay for user experience
// reasons.
func SoftReset() {
	time.Sleep(200 * time.Millisecond)
	machine.CPUReset()
}
