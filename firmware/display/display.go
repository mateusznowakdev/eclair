package display

import (
	"device/sam"
	"machine"

	"tinygo.org/x/drivers/ssd1306"

	"eclair/peripherals"
)

const displayWidth = 128
const displayHeight = 32

const ContrastHigh = 255
const ContrastNormal = 63
const ContrastLow = 0

type Display struct {
	device ssd1306.Device

	contrast uint8
	inverted bool
}

func newSPI() machine.SPI {
	frequency := uint32(peripherals.PatchedGCLK0Frequency(8_000_000))

	spi := machine.SPI{Bus: sam.SERCOM0_SPI, SERCOM: 0}
	_ = spi.Configure(machine.SPIConfig{SCK: machine.PA09, SDO: machine.PA11, SDI: machine.PA08, Frequency: frequency})

	return spi
}

// NewDisplay creates a new Display instance, configures a SERCOM bus in SPI
// mode, and sends the initialization sequence to the display.
func NewDisplay() Display {
	spi := newSPI()

	device := ssd1306.NewSPI(&spi /* DC: */, machine.PA07 /* RST: */, machine.PA06 /* CS: */, machine.PA05)
	device.Configure(ssd1306.Config{Width: displayWidth, Height: displayHeight})

	// reduce device brightness and allow for wider brightness range
	device.Command(ssd1306.SETPRECHARGE)
	device.Command(0)
	device.Command(ssd1306.SETVCOMDETECT)
	device.Command(0)

	display := Display{device: device}
	display.SetContrast(ContrastNormal)
	display.SetInverted(false)

	return display
}

// ClearBuffer clears the display buffer.
func (d *Display) ClearBuffer() {
	d.device.ClearBuffer()
}

// ClearDisplay clears the display buffer and sends it to the display.
func (d *Display) ClearDisplay() {
	d.device.ClearDisplay()
}

// Contrast returns the current contrast value.
func (d *Display) Contrast() uint8 {
	return d.contrast
}

// Display sends the buffer data to the display.
func (d *Display) Display() {
	_ = d.device.Display()
}

// Inverted returns the current invert state.
func (d *Display) Inverted() bool {
	return d.inverted
}

// SetContrast changes the contrast value and stores it in the Display instance
// (because the display is in write-only mode).
func (d *Display) SetContrast(value uint8) {
	d.device.Command(ssd1306.SETCONTRAST)
	d.device.Command(value)
	d.contrast = value
}

// SetInverted changes the invert state and stores it in the Display instance
// (because the display is in write-only mode).
func (d *Display) SetInverted(value bool) {
	if value {
		d.device.Command(ssd1306.INVERTDISPLAY)
	} else {
		d.device.Command(ssd1306.NORMALDISPLAY)
	}

	d.inverted = value
}
