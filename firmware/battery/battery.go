package battery

import "machine"

const threshold = 3400 // mV
const hysteresis = 100 // mV

type Battery struct {
	adc     machine.ADC
	voltage int
	good    bool
}

// NewBattery creates a new Battery instance and configures the underlying ADC
// channel.
func NewBattery() Battery {
	machine.InitADC()

	adc := machine.ADC{Pin: machine.PA02}
	adc.Configure(machine.ADCConfig{})

	return Battery{adc: adc}
}

func (b *Battery) refresh() {
	// the maximum raw value is 65535 and maximum read voltage is 6600mV
	// dividing raw value by 10 seems to give more accurate results
	b.voltage = int(b.adc.Get()) / 10

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
