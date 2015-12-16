package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/howeyc/fsnotify"
	// "github.com/k0kubun/pp"
)

type project struct {
	dirs      []string
	goFiles   []string
	jadeFiles []string
	makeLock  bool
	cmd       *exec.Cmd
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

			if filepath.Ext(fname) == ".go" {
				dot.goFiles = append(dot.goFiles, wd+"/"+fname)
				continue
			}
			if filepath.Ext(fname) == ".jade" {
				dot.jadeFiles = append(dot.jadeFiles, wd+"/"+fname)
				continue
			}

			if file.IsDir() == true && fname[0] != '.' {
				dot.dirs = append(dot.dirs, wd+"/"+fname)
				continue
			}
		}
	}
}

func (dot *project) make() {
	mk := exec.Command("go", "build", "-o", "a.out")
	mk.Stdout = os.Stdout
	mk.Stderr = os.Stderr
	err := mk.Run()
	if err != nil {
		log.Println("Fail run make() - ", err)
		os.Exit(2)
	}
}

func (dot *project) start() {
	dot.cmd = exec.Command("a.out")
	dot.cmd.Stdout = os.Stdout
	dot.cmd.Stderr = os.Stderr

	go dot.cmd.Run()
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
		var event string
		for {
			select {
			case ev := <-watcher.Event:
				if event == fmt.Sprint(time.Now().Unix(), ev) {
					break
				}
				event = fmt.Sprint(time.Now().Unix(), ev)
				log.Println("-- make")
				pro.make()
				log.Println("-- stop")
				pro.stop()
				log.Println("-- start")
				pro.start()

			case err := <-watcher.Error:
				log.Println("error:", err)
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

	// pro.make()

	quit := make(chan bool)
	for {
		select {
		case <-quit:
			return
		}
	}
}
