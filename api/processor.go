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
	IsRunning() bool
	SupportedProperties() []Property
	SupportedRelationships() []Relationship
}

type BaseProcessor struct {
	name                   string
	id                     string
	isRunning              bool
	supportedProperties    []Property
	supportedRelationships []Relationship
}

func NewBaseProcessor(name string, id string) *BaseProcessor {
	return &BaseProcessor{name: name, id: id}
}

func (b *BaseProcessor) Name() string {
	return b.name
}

func (b *BaseProcessor) Id() string {
	return b.id
}

func (b *BaseProcessor) IsRunning() bool {
	return b.isRunning
}

func (b *BaseProcessor) SupportedProperties() []Property {
	return b.supportedProperties
}

func (b *BaseProcessor) SupportedRelationships() []Relationship {
	return b.supportedRelationships
}
