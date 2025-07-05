package display

type Line struct {
	Start int
	End   int
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
