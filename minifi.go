package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"time"
)

const (
	Day = time.Hour * 24
)

var propertyPattern = regexp.MustCompile(`(?m)^([^#].*)=((.*\\\r?\n)*.*)$`)

type Property interface {
	Name() string
	Description() string
	Value() string
}

func main() {
	nifiProps, err := ioutil.ReadFile("/Users/brosander/Github/nifi/nifi-assemblies/nifi-assembly/target/nifi-1.1.0-SNAPSHOT-bin/nifi-1.1.0-SNAPSHOT/conf/nifi.properties")
	if err != nil {
		panic(err.Error())
	}
	for _, match := range propertyPattern.FindAllSubmatch(nifiProps, -1) {
		fmt.Println(string(match[1]), ": ", string(match[2]))
	}
}
