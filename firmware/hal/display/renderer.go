package display

const cursorData = 0xAAAA

const maxSpriteLines = displayHeight / 8
const maxTextLines = maxSpriteLines - 1

type Line struct {
	Start int
	End   int
}

func drawCursor(buffer []byte, bufferPos uint) {
	buffer[bufferPos] = cursorData & 0xFF
	buffer[bufferPos+displayWidth] = cursorData >> 8
}

func drawSprite(sprite []uint8, buffer []byte, line uint, xOffset uint) {
	if line >= maxSpriteLines {
		return
	}

	bufferStart := line * displayWidth

	for _, col := range sprite {
		if xOffset >= displayWidth {
			break
		}

		buffer[bufferStart+xOffset] = col
		xOffset += 1
	}
}

func drawText(text []byte, cursor int, buffer []byte, line uint, xOffset uint) {
	if line >= maxTextLines {
		return
	}

	bufferStart := line * displayWidth

	for charNo, char := range text {
		glyph := getGlyph(char)

		for colNo, col := range glyph {
			if xOffset >= displayWidth {
				break
			}

			if charNo == cursor && colNo == 0 {
				col |= cursorData
			}

			buffer[bufferStart+xOffset] = uint8(col >> 0)
			buffer[bufferStart+xOffset+displayWidth] = uint8(col >> 8)

			xOffset += 1
		}
	}

	if cursor == len(text) {
		if len(text) > 0 {
			xOffset -= 1
		}
		drawCursor(buffer, bufferStart+xOffset)
	}
}

func getGlyph(char byte) []uint16 {
	char = char & 0x7F

	if char < 32 {
		return font[len(font)-1]
	}

	glyph := font[char-32]
	if glyph == nil {
		return font[len(font)-1]
	}

	return glyph
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

	if ll >= 1 {
		line := lines[ll-1]
		drawText(text[line.Start:line.End], cursor-line.Start, d.device.GetBuffer(), 2, 0)

		if ll >= 2 {
			line = lines[ll-2]
			drawText(text[line.Start:line.End], -1, d.device.GetBuffer(), 0, 0)
		}
	}
}

// DrawSprite copies a sprite data to the display buffer at given position.
//
// The line is a value between 0 (top) and 3 (bottom of the screen).
func (d *Display) DrawSprite(sprite []uint8, line uint, xOffset uint) {
	drawSprite(sprite, d.device.GetBuffer(), line, xOffset)
}

// DrawText copies the defined text to the display buffer at given X offset.
//
// The line is a value between 0 (top) and 2 (bottom of the screen).
func (d *Display) DrawText(text []byte, line uint, xOffset uint) {
	drawText(text, -1, d.device.GetBuffer(), line, xOffset)
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
		glyphWidth := len(getGlyph(char))
		wordLen += glyphWidth

		// handle whitespace between words
		if text[charNo] == ' ' {
			newLineLen := lineLen + wordLen

			if brokeLongWord || newLineLen-glyphWidth > displayWidth {
				lineLen = wordLen
				brokeLongWord = false
				updateLines(wordStart, wordStart)
			} else if newLineLen > displayWidth {
				lineLen = 0
				updateLines(charNo, charNo+1)
			} else {
				lineLen += wordLen
			}

			wordStart = charNo + 1
			wordLen = 0
		}

		// handle very long words
		if wordLen > displayWidth {
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
	if brokeLongWord || newLineLen > displayWidth {
		updateLines(wordStart, wordStart)
	}

	// add final line
	lines[len(lines)-1].End = len(text)

	return lines
}
