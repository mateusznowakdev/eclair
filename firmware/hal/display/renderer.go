package display

const spriteHeight = 8

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

const (
	MaskNone MaskType = iota
	MaskAll
	MaskCustom
)

type Alignment int

type MaskType int

type Line struct {
	Start int
	End   int
}

func drawCursor(buffer []byte, x int, y int) {
	drawSprite8(buffer, cursor, 0, x, y, MaskAll, nil)
	drawSprite8(buffer, cursor, 0, x, y+8, MaskAll, nil)
}

func drawSprite8(buffer []uint8, spritesheet [][]uint8, id int, x int, y int, mask MaskType, masksheet [][]uint8) {
	start := (y/8)*Width + x
	shift := y % 8

	for col := range spritesheet[id] {
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

		// It so happens that this is a very similar code as above,
		// but with OR instead of XOR+AND, and spriteData instead of maskData.

		colSprite := spritesheet[id][col]

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

func drawSprite16(buffer []uint8, spritesheet [][]uint8, id int, x int, y int, mask MaskType, masksheet [][]uint8) {
	count := len(spritesheet) / 2
	drawSprite8(buffer, spritesheet, id, x, y, mask, masksheet)
	drawSprite8(buffer, spritesheet, count+id, x, y+8, mask, masksheet)
}

func drawText(buffer []uint8, text []byte, x int, y int, align Alignment, cursor int) {
	sprites := make([]int, len(text))
	widths := make([]int, len(text))
	widthTotal := 0

	for charNo, char := range text {
		id, w := getGlyphIDWidth(char)
		sprites[charNo] = id
		widths[charNo] = w
		widthTotal += w
	}

	switch align {
	case AlignCenter:
		x -= widthTotal/2 + 1
	case AlignRight:
		x -= widthTotal + 1
	default:
	}

	for charNo, spriteId := range sprites {
		drawSprite16(buffer, font, spriteId, x, y, MaskAll, nil)
		if cursor == charNo {
			drawCursor(buffer, x, y)
		}

		x += widths[charNo]
	}

	if cursor == len(text) {
		drawCursor(buffer, x-1, y)
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

// DrawMultiText fills the entire display buffer with the defined text.
//
// This function splits the text into multiple lines that can fit horizontally
// on the screen, and then scrolls it so the cursor is always displayed on the
// bottom line.
//
// If cursor is set to a negative value, it is not displayed.
func (d *Display) DrawMultiText(text []byte, cursor int) {
	lines := GetLines(text)

	// remove lines until the cursor is on the bottom line
	for ln := len(lines) - 1; ln > 0; ln-- {
		if cursor < lines[ln].Start {
			lines = lines[:ln]
		}
	}

	// making code more readable by storing this in a variable
	ll := len(lines)

	if ll > 0 {
		line := lines[ll-1]
		drawText(d.device.GetBuffer(), text[line.Start:line.End], 0, 16, AlignLeft, cursor-line.Start)

		if ll > 1 {
			line = lines[ll-2]
			drawText(d.device.GetBuffer(), text[line.Start:line.End], 0, 0, AlignLeft, -1)
		}
	}
}

// DrawSprite8 copies sprite data (8px tall) to the display buffer at given
// position.
func (d *Display) DrawSprite8(spritesheet [][]uint8, id int, x int, y int, mask MaskType, masksheet [][]uint8) {
	drawSprite8(d.device.GetBuffer(), spritesheet, id, x, y, mask, masksheet)
}

// DrawSprite16 copies sprite data (16px tall) to the display buffer at given
// position.
func (d *Display) DrawSprite16(spritesheet [][]uint8, id int, x int, y int, mask MaskType, masksheet [][]uint8) {
	drawSprite16(d.device.GetBuffer(), spritesheet, id, x, y, mask, masksheet)
}

// DrawText renders text to the display buffer at given position.
func (d *Display) DrawText(text []byte, x int, y int, align Alignment) {
	drawText(d.device.GetBuffer(), text, x, y, align, -1)
}

// GetLines splits the defined text into multiple lines and returns the starting
// and ending points for each line. This function is exported so it can be used
// directly by apps, not only within the text renderer.
func GetLines(text []byte) []Line {
	// avoid re-allocations by creating larger backing array ahead of time
	lines := make([]Line, 0, len(text)/16)
	lines = append(lines, Line{Start: 0})

	// inner function, it can update the variable directly
	updateLines := func(end int, start int) {
		lines[len(lines)-1].End = end
		lines = append(lines, Line{Start: start})
	}

	lineLen := 0

	wordStart := 0
	wordLen := 0

	brokeLongWord := false

	for charNo, char := range text {
		_, glyphWidth := getGlyphIDWidth(char)
		wordLen += glyphWidth

		// handle whitespace between words
		if text[charNo] == ' ' {
			newLineLen := lineLen + wordLen

			if brokeLongWord || newLineLen-glyphWidth > Width {
				lineLen = wordLen
				brokeLongWord = false
				updateLines(wordStart, wordStart)
			} else if newLineLen > Width {
				lineLen = 0
				updateLines(charNo, charNo+1)
			} else {
				lineLen += wordLen
			}

			wordStart = charNo + 1
			wordLen = 0
		}

		// handle very long words
		if wordLen > Width {
			if wordStart > 0 {
				updateLines(wordStart, wordStart)
			}

			lineLen = glyphWidth

			wordStart = charNo
			wordLen = glyphWidth

			brokeLongWord = true
		}
	}

	// handle last word
	newLineLen := lineLen + wordLen
	if brokeLongWord || newLineLen > Width {
		updateLines(wordStart, wordStart)
	}

	// add final line
	lines[len(lines)-1].End = len(text)

	return lines
}
