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
		case 114: // "r":
			if _, err := os.Stat("./" + dot.name); err == nil {
				dot.stop()
				dot.start()
			} else {
				c.Infof("file '%s' does not exist", dot.name)
			}
		case 98: // "b":
			make_notify <- true
		case 106: // "j":
			dot.rebuildAllJade()
		case 32: // " ":
			if pause {
				go dot.make()
				pause = false
				c.Info("start make")
			} else {
				make_notify <- false
				pause = true
				c.Info("stop  make")
			}
		case 113: //"q":
			dot.stop()
			c.Info("goodbye")
			break Bye
		default:
			fmt.Println(b)
		}
	}
}
