//go:build sam && atsamd21 && eclair_m0

package machine

// used to reset into bootloader
const resetMagicValue = 0xf01669ef

// USBCDC pins
const (
	USBCDC_DM_PIN = PA24
	USBCDC_DP_PIN = PA25
)

// UART pins
const (
	UART_TX_PIN = NoPin
	UART_RX_PIN = NoPin
)

// SPI pins
const (
	SPI0_SCK_PIN = PA09
	SPI0_SDO_PIN = PA11
	SPI0_SDI_PIN = PA08
)

// SPI on the Eclair M0
var SPI0 = sercomSPIM0

// I2C pins
const (
	SDA_PIN = NoPin
	SCL_PIN = NoPin
)

// I2S pins
const (
	I2S_SCK_PIN = NoPin
	I2S_SDO_PIN = NoPin
	I2S_SDI_PIN = NoPin
	I2S_WS_PIN  = NoPin
)

// Display pins
const (
	DISP_DC_PIN  = PA07
	DISP_CS_PIN  = PA05
	DISP_RST_PIN = PA06
)

// Keypad pins
const (
	KEYS_COL1_PIN = PA23
	KEYS_COL2_PIN = PA22
	KEYS_COL3_PIN = PA19
	KEYS_COL4_PIN = PA18
	KEYS_COL5_PIN = PA01
	KEYS_ROW1_PIN = PA04
	KEYS_ROW2_PIN = PA28
	KEYS_ROW3_PIN = PA16
)

// Voltage measurement pins
const VMETER_PIN = PA02

// USB CDC identifiers
const (
	usb_STRING_PRODUCT      = "EclairM0"
	usb_STRING_MANUFACTURER = "Mateusz Nowak"
)

var (
	usb_VID uint16 = 0x1209
	usb_PID uint16 = 0xEC1A
)
