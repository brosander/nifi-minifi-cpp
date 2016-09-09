package main

import (
	"fmt"
	"io/ioutil"
)

type GetFile struct {
}

func (getFile *GetFile) OnTrigger(context *ProcessContext, session *ProcessSession) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
