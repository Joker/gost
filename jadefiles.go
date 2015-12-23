package main

import (
	"fmt"
	"io/ioutil"
	// "log"
	// "os"
	// "os/exec"

	c "github.com/Joker/csi"
	"github.com/Joker/jade"
)

func (dot *project) jade() {
	var fileName string
	for {
		fileName = <-jade_notify
		dat, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Printf("ReadFile error: %v", err)
			return
		}

		c.Info(fileName)
		tmpl, err := jade.Parse("jt", string(dat))
		if err != nil {
			fmt.Printf("ReadFile error: %v", err)
			return
		}

		fmt.Println(tmpl)
	}
}
