// +build linux
package main

import (
	"fmt"
	"os"
	"os/exec"
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
		switch string(b) {
		case "r":
			// TODO check ./a.out
			dot.stop()
			dot.start()
		case "b":
			make_notify <- true
		case " ":
			if pause {
				go dot.make()
				pause = false
				fmt.Println("start make")
			} else {
				make_notify <- false
				pause = true
				fmt.Println("stop  make")
			}
		case "q":
			dot.stop()
			fmt.Println("goodbye")
			break Bye
		}
	}
}
