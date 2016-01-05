// +build linux

package main

import (
	"fmt"
	"os"

	c "github.com/Joker/ioterm"
)

func (dot *project) command() {

	c.RawMode()
	defer c.OrigMode()

	var (
		b     = make([]byte, 1)
		pause bool
	)
Bye:
	for {
		os.Stdin.Read(b)

		switch uint32(b[0]) {
		case 104: // "h"
			commandHelp()
		case 114: // "r"
			if _, err := os.Stat("./" + dot.name); err == nil {
				dot.stop()
				go dot.start()
			} else {
				c.Infof("file '%s' does not exist", dot.name)
			}
		case 98: // "b"
			make_notify <- true
		case 106: // "j"
			dot.rebuildAllJade()
		case 32: // " "
			if pause {
				go dot.make()
				pause = false
				c.Info("start make")
			} else {
				make_notify <- false
				pause = true
				c.Info("stop  make")
			}
		case 10: // "enter"			ScrollUp 7 lines
			c.Esc("7S")
		case 127: // "backspace" 	Clean screen
			c.Esc("2J")
		case 113: // "q"
			dot.stop()
			c.Info("goodbye")
			break Bye
		default:
			fmt.Println(b)
		}
	}
}

func commandHelp() {
	fmt.Println(`
	q - quit
	h - print help
	j - compile all .jade files
	r - restart go program
	b - rebuild go program
	space - stop/start build go program
	enter - ScrollUp 7 lines
	backspace - Clean screen
	`)
}
