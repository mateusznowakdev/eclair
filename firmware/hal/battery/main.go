package battery

import (
	"device/sam"
	"machine"
)

const threshold = 3300 // mV
const hysteresis = 100 // mV

type Battery struct {
	adc     machine.ADC
	voltage int
	good    bool
}

// New creates a new Battery instance and configures the underlying ADC channel
// using 1V fixed reference voltage.
func New() Battery {
	machine.InitADC()

	adc := machine.ADC{Pin: machine.VMETER_PIN}
	adc.Configure(machine.ADCConfig{})

	sam.ADC.INPUTCTRL.ReplaceBits(sam.ADC_INPUTCTRL_GAIN_1X, 0xF, sam.ADC_INPUTCTRL_GAIN_Pos)
	sam.ADC.REFCTRL.ReplaceBits(sam.ADC_REFCTRL_REFSEL_INT1V, 0xF, sam.ADC_REFCTRL_REFSEL_Pos)

	return Battery{adc: adc}
}

func (b *Battery) refresh() {
	b.voltage = int(b.adc.Get()) * 5250 / 65535

	if b.voltage < threshold {
		b.good = false
	} else if b.voltage > threshold+hysteresis {
		b.good = true
	}
}

// Good reads the current battery voltage, compares it to the previous value
// with hysteresis, and returns the boolean value based on a comparison result.
func (b *Battery) Good() bool {
	b.refresh()
	return b.good
}

// Voltage reads the current battery voltage and returns the value in
// millivolts.
func (b *Battery) Voltage() int {
	b.refresh()
	return b.voltage
}
