package api

type Property struct {
	Name         string
	Description  string
	DefaultValue string
}

type ProcessContext interface {
}

type ProcessSession interface {
	Create() FlowFile
	ImportFrom(fileName string, keepSourceFile bool, destination FlowFile) FlowFile
	PutAllAttributes(flowFile FlowFile, attributes map[string]string) FlowFile
	Transfer(flowFile FlowFile, relationship Relationship)
}

type Processor interface {
	Name() string
	Id() string
	GetProperty(name string) string
	IsRunning() bool
}
