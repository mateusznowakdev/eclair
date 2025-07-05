package launcher

import (
	"eclair/apps/erase"
	"eclair/apps/flashlight"
	"eclair/apps/mouse"
	"eclair/apps/notes"
)

var apps = []*Entry{
	{0, func() { notes.Run(notes.DefaultName) }},
	{1, func() { notes.Run("note2.txt") }},
	{2, func() { notes.Run("note3.txt") }},
	{3, func() { notes.Run("note4.txt") }},
	{4, func() { notes.Run("note5.txt") }},
	{5, mouse.Run},
	{6, flashlight.Run},
	nil,
	nil,
	{7, erase.Run},
}
