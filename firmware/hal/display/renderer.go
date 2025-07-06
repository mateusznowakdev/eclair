package display

const spriteHeight = 8

const (
	AlignLeft AlignType = iota
	AlignCenter
	AlignRight
)

const (
	MaskNone MaskType = iota
	MaskAll
	MaskCustom
)

type AlignType int
type MaskType int

func drawCursor(buffer []byte, x int, y int) {
	drawSprite8(buffer, cursor, 0, x, y, AlignLeft, MaskAll, nil)
	drawSprite8(buffer, cursor, 0, x, y+8, AlignLeft, MaskAll, nil)
}

func drawSprite8(buffer []uint8, spritesheet [][]uint8, id int, x int, y int, align AlignType, mask MaskType, masksheet [][]uint8) {
	sprite := spritesheet[id]
	spriteWidth := len(sprite)

	switch align {
	case AlignCenter:
		x -= spriteWidth / 2
	case AlignRight:
		x -= spriteWidth
	default:
	}

	if x <= -spriteWidth || x >= Width || y <= -spriteHeight || y >= Height {
		return
	}

	start := (y/8)*Width + x
	shift := y % 8

	for col := range sprite {
		if x < 0 || x >= Width {
			x++
			continue
		}

		if mask != MaskNone {
			var colMask uint8

			switch mask {
			case MaskAll:
				colMask = 0xFF
			case MaskCustom:
				colMask = masksheet[id][col]
			default:
			}

			if shift == 0 {
				// fast aligned to the page
				buffer[start+col] &= 0xFF ^ colMask
			} else if y < 0 {
				// only the bottom pixels
				buffer[start+col] &= 0xFF ^ (colMask >> -shift)
			} else if y >= Height-spriteHeight {
				// only the top pixels
				buffer[start+col] &= 0xFF ^ (colMask << shift)
			} else {
				// full spriteData
				buffer[start+col] &= 0xFF ^ (colMask << shift)
				buffer[start+col+Width] &= 0xFF ^ (colMask >> (spriteHeight - shift))
			}
		}

		// it so happens that the code below is very similar to the one above,
		// but with OR instead of XOR+AND, and spriteData instead of maskData

		colSprite := sprite[col]

		if shift == 0 {
			buffer[start+col] |= colSprite
		} else if y < 0 {
			buffer[start+col] |= colSprite >> -shift
		} else if y >= Height-spriteHeight {
			buffer[start+col] |= colSprite << shift
		} else {
			buffer[start+col] |= colSprite << shift
			buffer[start+col+Width] |= colSprite >> (spriteHeight - shift)
		}

		x++
	}
}

func drawSprite16(buffer []uint8, spritesheet [][]uint8, id int, x int, y int, align AlignType, mask MaskType, masksheet [][]uint8) {
	count := len(spritesheet) / 2
	drawSprite8(buffer, spritesheet, id, x, y, align, mask, masksheet)
	drawSprite8(buffer, spritesheet, count+id, x, y+8, align, mask, masksheet)
}

func drawText(buffer []uint8, text []byte, x int, y int, align AlignType, cursor int) {
	sprites := make([]int, len(text))
	widths := make([]int, len(text))
	widthTotal := 0

	for charNo, char := range text {
		id, w := getGlyphIDWidth(char)
		sprites[charNo] = id
		widths[charNo] = w
		widthTotal += w
	}

	// glyphs have more padding to the left than to the right,
	// therefore this code is slightly different to drawSprite8
	switch align {
	case AlignCenter:
		x -= widthTotal/2 + 1
	case AlignRight:
		x -= widthTotal + 1
	default:
	}

	for charNo, spriteId := range sprites {
		drawSprite16(buffer, font, spriteId, x, y, AlignLeft, MaskAll, nil)
		if cursor == charNo {
			drawCursor(buffer, x, y)
		}

		x += widths[charNo]
	}

	if cursor == len(text) {
		if cursor > 0 {
			x--
		}
		drawCursor(buffer, x, y)
	}
}

func getGlyphIDWidth(char byte) (int, int) {
	first := 0x20
	missing := 0x7F

	idx := int(char) - first
	if idx < 0 || idx >= len(font) {
		idx = missing - first
	}

	return idx, len(font[idx])
}

// DrawSprite8 copies sprite data (8px tall) to the display buffer at given
// position.
func (d *Display) DrawSprite8(spritesheet [][]uint8, id int, x int, y int, align AlignType, mask MaskType, masksheet [][]uint8) {
	drawSprite8(d.device.GetBuffer(), spritesheet, id, x, y, align, mask, masksheet)
}

// DrawSprite16 copies sprite data (16px tall) to the display buffer at given
// position.
func (d *Display) DrawSprite16(spritesheet [][]uint8, id int, x int, y int, align AlignType, mask MaskType, masksheet [][]uint8) {
	drawSprite16(d.device.GetBuffer(), spritesheet, id, x, y, align, mask, masksheet)
}

// DrawText renders text to the display buffer at given position.
func (d *Display) DrawText(text []byte, x int, y int, align AlignType) {
	drawText(d.device.GetBuffer(), text, x, y, align, -1)
}
