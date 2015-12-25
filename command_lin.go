// +build linux
package main

import (
	"fmt"
	"os"
	"os/exec"

	c "github.com/Joker/csi"
)

func (dot *project) command() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var (
		b     []byte = make([]byte, 1)
		pause bool
	)
Bye:
	for {
		os.Stdin.Read(b)

		switch uint32(b[0]) {
		case 114: // "r":
			// TODO check ./a.out
			dot.stop()
			dot.start()
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
