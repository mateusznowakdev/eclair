package launcher

import (
	"eclair/apps/bootloader"
	"eclair/apps/flashlight"
	"eclair/apps/mouse"
	"eclair/apps/notes"
)

var apps = []Entry{
	{"notes", notes.Run},
	{"mouse", mouse.Run},
	{"flashlight", flashlight.Run},
	{"bootloader", bootloader.Run},
}
