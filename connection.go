package main

type Connection interface {
	Name() string
	Id() string
	SourceProcessor() Processor
	DestinationProcessor() Processor
	SourceRelationship() Relationship
	MaxQueueSize() uint64
	FlowExpirationDuration() uint64
	IsEmpty() bool
	QueueDataSize() uint64
	Put(flowFileRecord string)
	Drain()
}
