package notes

import (
	"errors"
	"machine"
	"os"
	"slices"
	"time"

	"tinygo.org/x/tinyfs/littlefs"
)

const name = "note.txt"
const maxSize = 256
const writeDelay = 1500 // ms

type Note struct {
	fs *littlefs.LFS

	data      []byte
	cursor    int
	timestamp int64
}

func NewNote() (*Note, error) {
	fs := littlefs.New(machine.Flash)
	fs.Configure(&littlefs.Config{
		CacheSize:     256,
		LookaheadSize: 64,
		BlockCycles:   512,
	})

	err := fs.Mount()
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
	if len(n.data) < maxSize {
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
	f, err := n.fs.OpenFile(name, os.O_CREATE|os.O_RDONLY)
	if err != nil {
		return err
	}

	inf, _ := n.fs.Stat(name) // f.Stat().Size() is crashing
	size := inf.Size()
	if size > maxSize {
		return errors.New("file too large")
	}

	n.data = make([]byte, size, maxSize)
	n.cursor = int(size)

	if size > 0 {
		_, err = f.Read(n.data)
		if err != nil {
			return err
		}
	}

	_ = f.Close()
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

	f, err := n.fs.OpenFile("note.tmp", os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return err
	}

	_, err = f.Write(n.data)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	_ = n.fs.Rename("note.tmp", name) // atomic on littlefs

	return nil
}

func (n *Note) writeDelayed() (bool, error) {
	if n.timestamp == 0 || n.timestamp+writeDelay > time.Now().UnixMilli() {
		return false, nil
	}

	err := n.write()
	return err == nil, err
}
