package main

import (
	"io/ioutil"
	"time"

	c "github.com/Joker/ioterm"
	"github.com/Joker/jade"
	"strings"
)

func (dot *project) jade() {
	var fileName string
	for {
		fileName = <-jade_notify

		time.Sleep(500 * time.Millisecond)
		for len(jade_notify) > 0 {
			<-jade_notify
		}

		parse(fileName)
	}
}

func (dot *project) rebuildAllJade() {
	for _, fname := range dot.jadeFiles {
		parse(fname)
	}
}

func parse(fileName string) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		c.Errorf("ReadFile: %v", err)
		return
	}

	c.Info("parse: ", fileName)
	tmpl, err := jade.Parse("jt", string(dat))
	if err != nil {
		c.Errorf("Jade template: %v", err)
	}

	err = ioutil.WriteFile(strings.Replace(fileName, ".jade", ".tpl", 1), []byte(tmpl), 0644)
	if err != nil {
		c.Errorf("WriteFile: %v", err)
	}
	c.Note("\n\n", tmpl)
}
