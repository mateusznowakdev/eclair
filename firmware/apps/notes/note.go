package notes

import (
	"machine"
	"slices"
	"time"

	"github.com/mateusznowakdev/minifs"
)

const name = "note.txt"
const writeDelay = 1500 // ms

type Note struct {
	fs *minifs.Filesystem

	data      []byte
	cursor    int
	timestamp int64
}

func NewNote() (*Note, error) {
	fs, err := minifs.Configure(machine.Flash)
	if err != nil {
		return nil, err
	}

	return &Note{fs: fs}, nil
}

func (n *Note) cursorLeft() {
	if n.cursor > 0 {
		n.cursor -= 1
	}
}

func (n *Note) cursorRight() {
	if n.cursor < len(n.data) {
		n.cursor += 1
	}
}

func (n *Note) delete() {
	if n.cursor > 0 {
		n.data = slices.Delete(n.data, n.cursor-1, n.cursor)
		n.cursor -= 1
	}

	n.markDirty()
}

func (n *Note) dirty() bool {
	return n.timestamp > 0
}

func (n *Note) insert(b byte) {
	if len(n.data) < n.fs.MaxFileSize() {
		n.data = slices.Insert(n.data, n.cursor, b)
		n.cursor += 1
	}

	n.markDirty()
}

func (n *Note) last() byte {
	if n.cursor > 0 {
		return n.data[n.cursor-1]
	} else {
		return 0
	}
}

func (n *Note) markDirty() {
	n.timestamp = time.Now().UnixMilli()
}

func (n *Note) read() error {
	exists, err := n.fs.Exists(name)
	if err != nil {
		return err
	}

	if exists {
		data, err := n.fs.Read(name)
		if err != nil {
			return err
		}

		n.data = data
		n.cursor = len(n.data)
	} else {
		n.data = make([]byte, 0)
		n.cursor = 0
	}

	return nil
}

func (n *Note) replace(b byte) {
	if n.cursor > 0 {
		n.data[n.cursor-1] = b
	}

	n.markDirty()
}

func (n *Note) write() error {
	n.timestamp = 0

	err := n.fs.Write(name, n.data)
	if err != nil {
		return err
	}

	return nil
}

func (n *Note) writeDelayed() (bool, error) {
	if n.timestamp == 0 || n.timestamp+writeDelay > time.Now().UnixMilli() {
		return false, nil
	}

	err := n.write()
	return err == nil, err
}
