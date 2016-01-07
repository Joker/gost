package main

import (
	"fmt"
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
	name      string
	dirs      []string
	jadeFiles []string
	cmd       *exec.Cmd
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
			ext := filepath.Ext(ev.Name)
			switch {
			case ext == ".go" && conf.goFilesPause:
				log.Println("-- watcher.Event: ", ev)
				make_notify <- true
			case ext == ".jade" && conf.jadeFilesPause:
				log.Println("-- watcher.Event: ", ev)
				jade_notify <- ev.Name
			case ext == ".sass" && conf.sassFilesPause:
			case ext == ".sql" && conf.sqlFilesPause:
			case ext == ".js" && conf.jsFilesPause:
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

	dot.name = filepath.Base(wd)

	filepath.Walk(wd, func(path string, file os.FileInfo, err error) error {
		if filepath.Ext(file.Name()) == ".jade" {
			dot.jadeFiles = append(dot.jadeFiles, path)
		}
		if file.IsDir() == true {
			dot.dirs = append(dot.dirs, path)
		}
		return nil
	})

	// pp.Println(dot)

	return dot
}

func main() {
	pro := initProject()

	go pro.watch()
	go pro.make()
	go pro.jade()

	fmt.Println(">>>")
	pro.command()
}
