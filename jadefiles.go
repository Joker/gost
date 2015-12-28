package main

import (
	"fmt"
	"io/ioutil"
	"time"

	c "github.com/Joker/ioterm"
	"github.com/Joker/jade"
)

func (dot *project) jade() {
	var fileName string
	for {
		fileName = <-jade_notify

		time.Sleep(500 * time.Millisecond)
		for len(jade_notify) > 0 {
			<-jade_notify
		}
		dat, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Printf("ReadFile error: %v", err)
			return
		}

		c.Info(fileName)
		tmpl, err := jade.Parse("jt", string(dat))
		if err != nil {
			fmt.Printf("Jade template error: %v", err)
		}

		fmt.Println(tmpl)
	}
}

func (dot *project) rebuildAllJade() {
	for _, fname := range dot.jadeFiles {
		dat, err := ioutil.ReadFile(fname)
		if err != nil {
			fmt.Printf("ReadFile error: %v", err)
			return
		}

		tmpl, err := jade.Parse("jt", string(dat))
		if err != nil {
			c.Errorf("Jade template error: %v", err)
		}
		c.Info(fname, "\n\n", tmpl)
	}
}
