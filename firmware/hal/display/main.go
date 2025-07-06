package display

import (
	"encoding/base64"
	"machine"

	"eclair/hal/clocks"

	"tinygo.org/x/drivers/ssd1306"
)

const Width = 128
const Height = 32

const ContrastHigh = 255
const ContrastNormal = 63
const ContrastLow = 0

type Display struct {
	device ssd1306.Device

	contrast uint8
	inverted bool
}

// Configure creates a new Display instance, configures a SERCOM bus in SPI
// mode, and sends the initialization sequence to the display.
func Configure() *Display {
	frequency := uint32(clocks.PatchedGCLK0Frequency(8_000_000))

	spi := machine.SPI0
	_ = spi.Configure(machine.SPIConfig{Frequency: frequency})

	device := ssd1306.NewSPI(spi, machine.DISP_DC_PIN, machine.DISP_RST_PIN, machine.DISP_CS_PIN)
	device.Configure(ssd1306.Config{Width: Width, Height: Height})

	// reduce device brightness and allow for wider brightness range
	device.Command(ssd1306.SETPRECHARGE)
	device.Command(0)
	device.Command(ssd1306.SETVCOMDETECT)
	device.Command(0)

	display := Display{device: device}
	display.SetContrast(ContrastNormal)
	display.SetInverted(false)

	return &display
}

// ClearBuffer clears the display buffer.
func (d *Display) ClearBuffer() {
	d.device.ClearBuffer()
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

// Screenshot copies the buffer data onto serial console using base64 encoding.
// This can be useful for debugging (use the parse_base64_screenshot.py script).
func (d *Display) Screenshot() {
	buffer := d.device.GetBuffer()
	encoded := base64.StdEncoding.EncodeToString(buffer)

	print("Buffer data:\r\n")
	for start := 0; start < len(encoded); start += 80 {
		end := min(start+80, len(encoded))
		print(encoded[start:end], "\r\n")
	}
	print("Done.\r\n")
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
