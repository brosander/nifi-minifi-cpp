package getfile

import (
	"os"
	"testing"
	"time"

	"github.com/apache/gominifi/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func MockReadDir(t *testing.T, outputs map[string][]os.FileInfo) func(string) ([]os.FileInfo, error) {
	return func(dir string) ([]os.FileInfo, error) {
		result, ok := outputs[dir]
		assert.True(t, ok)
		return result, nil
	}
}

type MockFileInfo struct {
	mock.Mock
}

func (m *MockFileInfo) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockFileInfo) Size() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockFileInfo) Mode() os.FileMode {
	args := m.Called()
	return args.Get(0).(os.FileMode)
}

func (m *MockFileInfo) ModTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockFileInfo) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockFileInfo) Sys() interface{} {
	args := m.Called()
	return args.Get(0)
}

func TestAverage(t *testing.T) {
	getFile := NewGetFile("testGetFile", "abcd")
	mockFile := new(MockFileInfo)
	mockFile.On("IsDir").Return(false)
	mockFile.On("Name").Return("testFile")
	getFile.readDir = MockReadDir(t, map[string][]os.FileInfo{
		"/tmp/mgo": []os.FileInfo{mockFile},
	})
	testProcessContext := new(mocks.ProcessContext)
	testProcessSession := new(mocks.ProcessSession)
	flowFile := new(mocks.FlowFile)
	testProcessSession.On("Create").Return(flowFile)
	testProcessSession.On("PutAllAttributes", flowFile, mock.AnythingOfType("map[string]string")).Return(flowFile)
	testProcessSession.On("ImportFrom", "/tmp/mgo/testFile", false, flowFile).Return(flowFile)
	testProcessSession.On("Transfer", flowFile, success)
	getFile.OnTrigger(testProcessContext, testProcessSession)

	testProcessContext.AssertExpectations(t)
	testProcessSession.AssertExpectations(t)
	mockFile.AssertExpectations(t)
}
