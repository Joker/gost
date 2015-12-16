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
		<-make_notify
		log.Println("-- stop")
		dot.stop()
		log.Println("-- make")

		for i := len(make_notify); i > 0; i-- {
			fmt.Println("======== ", len(make_notify))
			<-make_notify
		}

		build := exec.Command("go", "build", "-o", "a.out")
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr
		err := build.Run()
		if err != nil {
			log.Printf("Command finished with error: %v \n", err)
			continue
		}

		log.Println("-- start")
		dot.start()
	}
}

func (dot *project) start() {
	dot.cmd = exec.Command("./a.out")
	dot.cmd.Stdout = os.Stdout
	dot.cmd.Stderr = os.Stderr
	go dot.cmd.Run()
}

func (dot *project) stop() {
	if dot.cmd != nil && dot.cmd.Process != nil {
		if !dot.cmd.ProcessState.Exited() {
			err := dot.cmd.Process.Kill()
			if err != nil {
				log.Println("Error cmd.Process.Kill() - ", err)
			}
		}
	}

	if rec := recover(); rec != nil {
		log.Println("Recovered in - ", rec)
	}
}

func main() {
	pro := project{}

	wd, err := os.Getwd()
	if err != nil {
		log.Println("Fail Getwd() - ", err)
		os.Exit(2)
	}
	pro.dirs = append(pro.dirs, wd)
	// pro.dirs = append(pro.dirs, "../test")
	// os.Chdir(path)

	var last, cursor int
	for {
		last = len(pro.dirs)
		pro.parseDir(cursor)
		if last == cursor {
			break
		}
		cursor = last
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Fail fsnotify NewWatcher() - ", err)
		os.Exit(2)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if filepath.Ext(ev.Name) == ".go" && ev.IsModify() {
					log.Println("-- watcher.Event: ", ev)
					make_notify <- true
				}
			case err := <-watcher.Error:
				log.Println("watcher.Error: ", err)
			}
		}
	}()

	for _, dir := range pro.dirs {
		err = watcher.Watch(dir)
		if err != nil {
			log.Println("Fail fsnotify Watch() - ", err)
			os.Exit(2)
		}
	}

	go pro.make()

	quit := make(chan bool)
	for {
		select {
		case <-quit:
			return
		}
	}
}
