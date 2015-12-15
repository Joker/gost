package main

import (
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
	"github.com/k0kubun/pp"
)


type workDir 	struct {
	dirs 		[]string
	goFiles 	[]string
	jadeFiles 	[]string
}	


func (dot *workDir) parseDir(sl int) {
	for _, wd := range dot.dirs[sl:] {
		fileInfo, err := ioutil.ReadDir(wd)
		if err != nil {
			log.Println("Fail ReadDir() - %s", err)
			os.Exit(2)
		}

		var fname string
		for _, file := range fileInfo {
			fname = file.Name()

			if filepath.Ext(fname) == ".rs" {
				dot.goFiles = append(dot.goFiles, wd+"/"+fname)
				continue
			}
			if filepath.Ext(fname) == ".yml" {
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


func main() {
	workFolder := workDir{}

	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Println("Fail Getwd() - %s", err)
	// 	os.Exit(2)
	// }
	// workFolder.dirs = append(workFolder.dirs, wd)
	workFolder.dirs = append(workFolder.dirs, "../git2-rs")

	var last, cursor int
	for {
		last = len(workFolder.dirs)
		workFolder.parseDir(cursor)
		if last == cursor { break }
		cursor = last
	}
	pp.Println(workFolder)
}