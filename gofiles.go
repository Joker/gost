package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	l "github.com/Joker/ioterm"
	c "github.com/Joker/ioterm/color"
	"github.com/Joker/panicparse/inter"
	// "github.com/k0kubun/pp"
)

var realPanic = false

type panicWriter struct {
	w io.Writer
}

func (pw panicWriter) Write(p []byte) (int, error) {
	pw.w.Write([]byte("\033[31m"))
	return pw.w.Write(p)
}

func (dot *project) make() {
	for {
		<-make_notify
		dot.stop()

		time.Sleep(time.Second * 1)
		for len(make_notify) > 0 {
			<-make_notify
		}

		fmt.Println(c.Blue_h, "--  make --", c.Reset)

		build := exec.Command("go", "build", "-o", dot.name)
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr
		err := build.Run()
		if err != nil {
			fmt.Printf("\n\n%s\n%s\nBuild finished with error: %v \n\n", c.Blue_b, c.Reset, err)
			continue
		}

		go dot.start()
	}
}

func (dot *project) start() {
	fmt.Println(c.Green_h, "-- start --", c.Reset, l.N(22), "\n", c.Reset)

	dot.cmd = exec.Command("./" + dot.name)
	dot.cmd.Stdout = os.Stdout
	// dot.cmd.Stderr = panicWriter{os.Stderr}
	stderr, err := dot.cmd.StderrPipe()
	if err != nil {
		log.Println("Error cmd.StderrPipe() - ", err)
	}

	err = dot.cmd.Start()
	if err != nil {
		log.Println("Error cmd.Start() - ", err)
	}

	sebuf := new(bytes.Buffer)
	go io.Copy(sebuf, stderr)

	err = dot.cmd.Wait()
	if err != nil {
		if sebuf.Len() > 0 {
			if conf.realPanic {
				fmt.Println(c.Red, sebuf, c.Reset)
			}
			pp, _ := internal.ParsePanic(sebuf)
			fmt.Println(string(pp))
		}
		l.Errorf("%s pid:%d -- %v", dot.name, err.(*exec.ExitError).Pid(), err)
	}
}

func (dot *project) stop() {
	if dot.cmd != nil && dot.cmd.Process != nil {
		fmt.Println(c.Red_h, "--  stop --", c.Reset)
		err := dot.cmd.Process.Kill()
		if err != nil {
			l.Note("cmd.Process.Kill() - ", err)
		}
	}

	if rec := recover(); rec != nil {
		log.Println("Recovered in - ", rec)
	}
}
