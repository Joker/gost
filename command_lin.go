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

	b := make([]byte, 1)

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
			conf.goFilesPause = !conf.goFilesPause
			if conf.goFilesPause {
				c.Info("start build go program")
			} else {
				c.Note("stop  build go program")
			}
		case 105: // "i"
			conf.jadeFilesPause = !conf.jadeFilesPause
			if conf.jadeFilesPause {
				c.Info("start compile .jade")
			} else {
				c.Note("stop  compile .jade")
			}
		case 112: // "p"
			conf.realPanic = !conf.realPanic
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
	i - stop/start compile .jade files
	p - show/hide real panic message
	r - restart go program
	b - rebuild go program
	space - stop/start build go program
	enter - ScrollUp 7 lines
	backspace - Clean screen
	`)
}
