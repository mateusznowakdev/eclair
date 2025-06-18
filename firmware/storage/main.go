package storage

import (
	"encoding/binary"
	"errors"
	"machine"
)

/*
Filesystem implementation, using CRUD-like API instead of the standard Go API.
Created because single write using littlefs at 8 MHz takes as much as 100ms.
This is 10-15x faster.

"Features":

  - No long names (8.3 instead)
  - Flat structure, no directories
  - File size is the same as erase block size (faster, prevents fragmentation)
  - Metadata size is the same as erase block size
  - Basic wear leveling
  - Basic atomicity (metadata updated after the data is written)
  - No extra redundancy features such as checksums
  - No timestamps, permissions, etc.

A valid metadata section starts with "SfsMetadataBlk**" text / "magic number"
where "*" is the number of currently stored files (up to 65536 because 2 bytes).
Then, each metadata entry is 16 bytes long:

  - 12 bytes of filename
  - 2  bytes of data size
  - 2  bytes of block number

On SAMD21 with 256 bytes erase block size, it's therefore possible to create up
to 15 files, 256 bytes each.
(TODO: maybe use two blocks to have more larger files? consider this AFTER
       the basic implementation is done)

Updating a file results in it being copied to the next available block (unless
it's already used), metadata block rewritten to the yet another next available
block, and then finally old metadata block erased.

If there are multiple valid metadata blocks (for example due to power loss in
the middle of operation), the first one is considered valid. If there are no
valid metadata blocks, then the filesystem is considered damaged.

The order of entries in the metadata block is not guaranteed.

Operations:

1. Exists
2. List
3. Read
4. Write
5. Rename
6. Delete
7. Format

Each function returns an error object, either from machine.Flash (read-write
failures) or from Filesystem implementation.
*/

const magic = "SfsMetadataBlk"

var errDamaged = errors.New("filesystem is damaged")

func blkSize() int64 {
	return machine.Flash.EraseBlockSize()
}

func blkToByte(blocks int64) int64 {
	return blocks * blkSize()
}

type File struct {
	size  uint16
	block uint16
}

type Filesystem struct {
	files map[string]File

	valid     bool
	metaBlock int64
}

func Configure() (*Filesystem, error) {
	// todo: search for and read metadata
	// todo: handle missing metadata, set "valid" respectively
	return nil, nil
}

func (fs *Filesystem) Delete(name string) error {
	if !fs.valid {
		return errDamaged
	}
	// todo: handle no file
	// todo: delete file reference from files list, write new metadata, erase old m/d

	return nil
}

func (fs *Filesystem) Exists(name string) (bool, error) {
	if !fs.valid {
		return false, errDamaged
	}

	_, ok := fs.files[name]
	return ok, nil
}

func (fs *Filesystem) Format() error {
	fs.files = make(map[string]File)
	fs.metaBlock = 0

	return fs.writeMetaBlock()
}

func (fs *Filesystem) List() ([]string, error) {
	if !fs.valid {
		return nil, errDamaged
	}

	names := make([]string, len(fs.files))
	fileNo := 0
	for name := range fs.files {
		names[fileNo] = name
	}

	return names, nil
}

func (fs *Filesystem) MaxSize() int {
	return int(blkSize())
}

func (fs *Filesystem) Read(name string) ([]byte, error) {
	if !fs.valid {
		return nil, errDamaged
	}
	// todo: handle no file
	// todo: handle file too large (invalid size info stored)
	// todo: read data from the block

	return nil, nil
}

func (fs *Filesystem) Rename(old string, new string) error {
	if !fs.valid {
		return errDamaged
	}
	// todo: handle no old file
	// todo: handle invalid new name (empty, too long)
	// todo: update files list, write new metadata, erase old m/d
	// todo: handle overwrite existing files (it's dict, should be handled automatically)

	return nil
}

func (fs *Filesystem) Write(name string, data []byte) error {
	if !fs.valid {
		return errDamaged
	}
	// todo: handle invalid name (empty, too long)
	// todo: handle too many files
	// todo: handle data too large
	// todo: write new data to the next available block
	// todo: write new metadata, erase old m/d

	return nil
}

func (fs *Filesystem) writeMetaBlock() error {
	// todo: find next suitable block for wear leveling, including solving collisions
	//       maybe store all used blocks, including metadata block, in some dict-like cache

	err := machine.Flash.EraseBlocks(fs.metaBlock, 1)
	if err != nil {
		return err
	}

	buf := make([]byte, blkSize())

	copy(buf, magic)
	binary.LittleEndian.PutUint16(buf[14:], uint16(len(fs.files)))

	entryNo := 1
	for name, file := range fs.files {
		pos := entryNo * 16
		copy(buf[pos:], name)
		binary.LittleEndian.PutUint16(buf[pos+12:], file.size)
		binary.LittleEndian.PutUint16(buf[pos+14:], file.block)
		entryNo++
	}

	_, err = machine.Flash.WriteAt(buf, blkToByte(fs.metaBlock))
	if err != nil {
		return err
	}

	fs.valid = true

	return nil
}

/* ----------------- old attempts

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
*/
