package launcher

import (
	"eclair/apps/bootloader"
	"eclair/apps/erase"
	"eclair/apps/flashlight"
	"eclair/apps/mouse"
	"eclair/apps/notes"
)

var apps = []Entry{
	{"note 1", func() { notes.Run(notes.DefaultName) }},
	{"note 2", func() { notes.Run("note2.txt") }},
	{"note 3", func() { notes.Run("note3.txt") }},
	{"mouse", mouse.Run},
	{"flashlight", flashlight.Run},
	{"format", erase.Run},
	{"bootloader", bootloader.Run},
}
