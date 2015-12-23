package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	c "github.com/Joker/csi"
)

func (dot *project) make() {
	for {
		if <-make_notify {
			dot.stop()

			for i := len(make_notify); i > 0; i-- {
				<-make_notify
			}

			fmt.Println(c.Blue_h, "--  make --", c.Reset)

			build := exec.Command("go", "build", "-o", "a.out")
			build.Stdout = os.Stdout
			build.Stderr = os.Stderr
			err := build.Run()
			if err != nil {
				fmt.Printf("\n\n%s\n%s\nBuild finished with error: %v \n\n", c.Blue_b, c.Reset, err)
				continue
			}

			dot.start()
		} else {
			break
		}
	}
}

func (dot *project) start() {
	fmt.Println(c.Green_h, "-- start --", c.Reset, c.N(22), "\n", c.Reset)

	dot.cmd = exec.Command("./a.out")
	dot.cmd.Stdout = os.Stdout
	dot.cmd.Stderr = os.Stderr
	dot.cmd.Start()
}

func (dot *project) stop() {
	fmt.Println(c.Red_h, "--  stop --", c.Reset)

	if dot.cmd != nil && dot.cmd.Process != nil {
		err := dot.cmd.Process.Kill()
		if err != nil {
			log.Println("Error cmd.Process.Kill() - ", err)
		}
	}

	if rec := recover(); rec != nil {
		log.Println("Recovered in - ", rec)
	}
}
