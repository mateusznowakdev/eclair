package storage

import (
	"bytes"
	"machine"
)

type Bounds struct {
	StartBlock uint
	EndBlock   uint
}

type File struct {
	Data []byte

	bounds Bounds
	block  uint
}

// NewFile creates a new File instance with defined filesystem boundaries,
// defaulting to the first possible block.
//
// It is necessary to call File.Read to read stored data.
func NewFile(bounds Bounds) *File {
	return &File{bounds: bounds, block: bounds.StartBlock}
}

// MaxSize returns the maximum allowed data size for a filesystem block,
// excluding the metadata section.
func (f *File) MaxSize() int {
	return int(machine.Flash.EraseBlockSize() - 4)
}

// Read attempts to read data stored in the flash memory, within defined
// filesystem boundaries.
//
// If a valid block was found, this function will return true.
//
// If there is a critical error, such as reading out of bounds, it will be
// returned.
func (f *File) Read() (bool, error) {
	blockSize := uint(machine.Flash.EraseBlockSize())

	data := make([]byte, blockSize)

	for b := f.bounds.StartBlock; b < f.bounds.EndBlock; b++ {
		_, err := machine.Flash.ReadAt(data, int64(b*blockSize))
		if err != nil {
			return false, err
		}

		if data[blockSize-4] != 'E' || data[blockSize-3] != 'F' || data[blockSize-2] != 0 {
			continue
		}

		size := min(int(data[blockSize-1]), f.MaxSize())

		f.Data = data[:size]
		f.block = b

		return true, nil
	}

	f.Data = make([]byte, 0, blockSize)
	f.block = f.bounds.StartBlock

	return false, nil
}

// Write attempts to write data stored in a File instance, within defined
// filesystem boundaries.
//
// Write tries to do perform a basic wear leveling by writing data to the next
// available flash block, and then it erases old block by overwriting the magic
// number with garbage data. In the worst case scenario, there would be two
// "valid" blocks, and most of the time the older one will take precedence.
//
// If there is a critical error, such as writing out of bounds, it will be
// returned.
func (f *File) Write() error {
	blockSize := uint(machine.Flash.EraseBlockSize())
	writeSize := uint(machine.Flash.WriteBlockSize())

	meta := bytes.Repeat([]byte{0xFF}, int(writeSize))
	meta[len(meta)-4] = 'E'
	meta[len(meta)-3] = 'F'
	meta[len(meta)-2] = 0
	meta[len(meta)-1] = byte(len(f.Data))

	block := f.block + 1
	if block >= f.bounds.EndBlock {
		block = f.bounds.StartBlock
	}

	err := machine.Flash.EraseBlocks(int64(block), 1)
	if err != nil {
		return err
	}
	_, err = machine.Flash.WriteAt(f.Data, int64(block*blockSize))
	if err != nil {
		return err
	}
	_, err = machine.Flash.WriteAt(meta, int64((block+1)*blockSize-writeSize))
	if err != nil {
		return err
	}

	if f.block != block {
		_, err = machine.Flash.WriteAt(make([]byte, writeSize), int64((f.block+1)*blockSize-writeSize))
		if err != nil {
			return err
		}
	}

	f.block = block

	return nil
}
