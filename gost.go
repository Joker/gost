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
var jade_notify = make(chan string, 100)

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
			if filepath.Ext(ev.Name) == ".jade" {
				log.Println("-- watcher.Event: ", ev)
				jade_notify <- ev.Name
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
	go pro.jade()

	fmt.Println(".")
	pro.command()
}
