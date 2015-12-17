package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/howeyc/fsnotify"
	// "github.com/k0kubun/pp"
)

var make_notify = make(chan bool, 100)

type project struct {
	dirs []string
	// goFiles   []string
	// jadeFiles []string
	cmd *exec.Cmd
}

func (dot *project) parseDir(sl int) {
	for _, wd := range dot.dirs[sl:] {
		fileInfo, err := ioutil.ReadDir(wd)
		if err != nil {
			log.Println("Fail ReadDir() - ", err)
			os.Exit(2)
		}

		var fname string
		for _, file := range fileInfo {
			fname = file.Name()

			// if filepath.Ext(fname) == ".go" {
			// 	dot.goFiles = append(dot.goFiles, wd+"/"+fname)
			// 	continue
			// }
			// if filepath.Ext(fname) == ".jade" {
			// 	dot.jadeFiles = append(dot.jadeFiles, wd+"/"+fname)
			// 	continue
			// }

			if file.IsDir() == true && fname[0] != '.' {
				dot.dirs = append(dot.dirs, wd+"/"+fname)
				continue
			}
		}
	}
}

func (dot *project) make() {
	for {
		if <-make_notify {
			log.Println("-- stop")
			dot.stop()

			for i := len(make_notify); i > 0; i-- {
				<-make_notify
			}

			log.Println("-- make\n")
			build := exec.Command("go", "build", "-o", "a.out")
			build.Stdout = os.Stdout
			build.Stderr = os.Stderr
			err := build.Run()
			if err != nil {
				fmt.Printf("\n\nBuild finished with error: %v \n\n", err)
				continue
			}

			log.Println("-- start\n\n")
			dot.start()
		} else {
			break
		}
	}
}

func (dot *project) start() {
	dot.cmd = exec.Command("./a.out")
	dot.cmd.Stdout = os.Stdout
	dot.cmd.Stderr = os.Stderr
	dot.cmd.Start()
}

func (dot *project) stop() {
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

func (dot *project) watch() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Fail fsnotify NewWatcher() - ", err)
		os.Exit(2)
	}

	for _, dir := range dot.dirs {
		err = watcher.Watch(dir)
		if err != nil {
			log.Println("Fail fsnotify Watch() - ", err)
			os.Exit(2)
		}
	}

	for {
		select {
		case ev := <-watcher.Event:
			if filepath.Ext(ev.Name) == ".go" {
				log.Println("-- watcher.Event: ", ev)
				make_notify <- true
			}
		case err := <-watcher.Error:
			log.Println("watcher.Error: ", err)
		}
	}
}

func initProject() project {
	dot := project{}

	wd, err := os.Getwd()
	if err != nil {
		log.Println("Fail Getwd() - ", err)
		os.Exit(2)
	}
	dot.dirs = append(dot.dirs, wd)

	var last, cursor int
	for {
		last = len(dot.dirs)
		dot.parseDir(cursor)
		if last == cursor {
			break
		}
		cursor = last
	}
	return dot
}

func main() {
	pro := initProject()

	go pro.watch()
	go pro.make()

	pro.command()
}
