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


func (dot *workDir) getSubdirectory() {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("Fail Getwd() - %s", err)
		os.Exit(2)
	}
	dot.dirs = append(dot.dirs, wd)

	fileInfo, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Println("Fail ReadDir() - %s", err)
		os.Exit(2)
	}

	var fname string
	for _, file := range fileInfo {
		fname = file.Name()

		if filepath.Ext(fname) == ".yml" {
			dot.goFiles = append(dot.goFiles, fname)
			continue
		}

		if file.IsDir() == true && fname[0] != '.' {
			dot.dirs = append(dot.dirs, wd+"/"+fname)
			continue
		}
	}
}

func main() {
	asd := workDir{}
	asd.getSubdirectory()
	pp.Println(asd)
}