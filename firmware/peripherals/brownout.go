package peripherals

import "device/sam"

// ConfigureBOD33 enables the brown-out detector, which resets the device on a
// significant voltage drop on a 3.3V bus, preventing system failures such as
// invalid flash memory operations.
func ConfigureBOD33() {
	// disable peripheral
	sam.SYSCTRL.BOD33.ClearBits(sam.SYSCTRL_BOD33_ENABLE)
	for !sam.SYSCTRL.PCLKSR.HasBits(sam.SYSCTRL_PCLKSR_B33SRDY) {
	}

	// set minimal system voltage to ~2.9V to match Adafruit uf2-samdx1 bootloader
	sam.SYSCTRL.BOD33.Set((39 << sam.SYSCTRL_BOD33_LEVEL_Pos) | (sam.SYSCTRL_BOD33_ACTION_RESET << sam.SYSCTRL_BOD33_ACTION_Pos) | sam.SYSCTRL_BOD33_HYST)

	// enable peripheral
	sam.SYSCTRL.BOD33.SetBits(sam.SYSCTRL_BOD33_ENABLE)
	for !sam.SYSCTRL.PCLKSR.HasBits(sam.SYSCTRL_PCLKSR_BOD33RDY) {
	}

	// wait for system voltage to be stable enough
	for sam.SYSCTRL.PCLKSR.HasBits(sam.SYSCTRL_PCLKSR_BOD33DET) {
	}
}
