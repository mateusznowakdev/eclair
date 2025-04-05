package bootloader

import "machine"

func Run() {
	machine.EnterBootloader()

	//goland:noinspection ALL
	for {
	}
}
