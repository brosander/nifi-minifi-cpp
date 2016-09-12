package main

import (
	"github.com/stretchr/testify/mock"
)

type MockProcessContext struct {
	mock.Mock
}

type MockProcessSession struct {
	mock.Mock
}

func (m *MockProcessSession) Create() FlowFile {
	args := m.Called()
	return args.Get(0).(FlowFile)
}

func (m *MockProcessSession) ImportFrom(file string, keepSourceFile bool, destination FlowFile) FlowFile {
	args := m.Called(file, keepSourceFile, destination)
	return args.Get(0).(FlowFile)
}

func (m *MockProcessSession) PutAllAttributes(flowFile FlowFile, attributes map[string]string) FlowFile {
	args := m.Called(flowFile, attributes)
	return args.Get(0).(FlowFile)
}

func (m *MockProcessSession) Transfer(flowFile FlowFile, relationship Relationship) {
	m.Called(flowFile, relationship)
}

type MockProcessor struct {
	mock.Mock
}

func (m *MockProcessor) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockProcessor) Id() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockProcessor) GetProperty(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func (m *MockProcessor) IsRunning() bool {
	args := m.Called()
	return args.Bool(0)
}
