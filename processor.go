package main

type ProcessContext interface {
}

type ProcessSession interface {
}

type Processor interface {
	Name() string
	Id() string
	GetProperty(name string) string
	IsRunning() bool
}
