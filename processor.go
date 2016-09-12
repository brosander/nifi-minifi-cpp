package main

import "os"

type ProcessContext interface {
}

type ProcessSession interface {
	Create() FlowFile
	ImportFrom(file *os.File, destination FlowFile) FlowFile
	PutAllAttributes(flowFile FlowFile, attributes map[string]string) FlowFile
}

type Processor interface {
	Name() string
	Id() string
	GetProperty(name string) string
	IsRunning() bool
}
