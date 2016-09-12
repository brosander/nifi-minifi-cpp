package mocks

import "github.com/apache/gominifi/api"
import "github.com/stretchr/testify/mock"

// ProcessSession is an autogenerated mock type for the ProcessSession type
type ProcessSession struct {
	mock.Mock
}

// Create provides a mock function with given fields:
func (_m *ProcessSession) Create() api.FlowFile {
	ret := _m.Called()

	var r0 api.FlowFile
	if rf, ok := ret.Get(0).(func() api.FlowFile); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.FlowFile)
		}
	}

	return r0
}

// ImportFrom provides a mock function with given fields: fileName, keepSourceFile, destination
func (_m *ProcessSession) ImportFrom(fileName string, keepSourceFile bool, destination api.FlowFile) api.FlowFile {
	ret := _m.Called(fileName, keepSourceFile, destination)

	var r0 api.FlowFile
	if rf, ok := ret.Get(0).(func(string, bool, api.FlowFile) api.FlowFile); ok {
		r0 = rf(fileName, keepSourceFile, destination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.FlowFile)
		}
	}

	return r0
}

// PutAllAttributes provides a mock function with given fields: flowFile, attributes
func (_m *ProcessSession) PutAllAttributes(flowFile api.FlowFile, attributes map[string]string) api.FlowFile {
	ret := _m.Called(flowFile, attributes)

	var r0 api.FlowFile
	if rf, ok := ret.Get(0).(func(api.FlowFile, map[string]string) api.FlowFile); ok {
		r0 = rf(flowFile, attributes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.FlowFile)
		}
	}

	return r0
}

// Transfer provides a mock function with given fields: flowFile, relationship
func (_m *ProcessSession) Transfer(flowFile api.FlowFile, relationship api.Relationship) {
	_m.Called(flowFile, relationship)
}
