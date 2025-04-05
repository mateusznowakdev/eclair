package notes

import (
	"slices"
	"time"

	"eclair/apps"
	"eclair/storage"
)

const writeDelay = 1500 // ms

type Note struct {
	file *storage.File

	cursor    int
	timestamp int64
}

func NewNote() Note {
	bounds := apps.Bounds["notes"]
	file := storage.NewFile(bounds)

	return Note{file: file}
}

func (n *Note) cursorLeft() {
	if n.cursor > 0 {
		n.cursor -= 1
	}
}

func (n *Note) cursorRight() {
	if n.cursor < len(n.file.Data) {
		n.cursor += 1
	}
}

func (n *Note) delete() {
	if n.cursor > 0 {
		n.file.Data = slices.Delete(n.file.Data, n.cursor-1, n.cursor)
		n.cursor -= 1
	}

	n.markDirty()
}

func (n *Note) dirty() bool {
	return n.timestamp > 0
}

func (n *Note) insert(b byte) {
	if len(n.file.Data) < n.file.MaxSize() {
		n.file.Data = slices.Insert(n.file.Data, n.cursor, b)
		n.cursor += 1
	}

	n.markDirty()
}

func (n *Note) last() byte {
	if n.cursor > 0 {
		return n.file.Data[n.cursor-1]
	} else {
		return 0
	}
}

func (n *Note) markDirty() {
	n.timestamp = time.Now().UnixMilli()
}

func (n *Note) read() (bool, error) {
	success, err := n.file.Read()
	n.cursor = len(n.file.Data)

	return success, err
}

func (n *Note) replace(b byte) {
	if n.cursor > 0 {
		n.file.Data[n.cursor-1] = b
	}

	n.markDirty()
}

func (n *Note) write() (bool, error) {
	n.timestamp = 0
	return true, n.file.Write()
}

func (n *Note) writeDelayed() (bool, error) {
	if n.timestamp == 0 || n.timestamp+writeDelay > time.Now().UnixMilli() {
		return false, nil
	}

	return n.write()
}
