package getfile

import (
	"os"
	"testing"

	apiMocks "github.com/apache/gominifi/api/mocks"
	"github.com/apache/gominifi/framework/processor"
	"github.com/apache/gominifi/processors/getfile/mocks"
	"github.com/stretchr/testify/mock"
)

func TestSingleFile(t *testing.T) {
	getFile := NewGetFile("testGetFile", "abcd")

	mockFile := new(mocks.MockFileInfo)
	mockFile.On("IsDir").Return(false)
	mockFile.On("Name").Return("testFile")

	getFile.readDir = mocks.MockReadDir(t, map[string][]os.FileInfo{
		"/tmp/mgo": []os.FileInfo{mockFile},
	})

	testProcessContext := new(apiMocks.ProcessContext)
	testProcessContext.On("GetPropertyValue", inputDirectory).Return(processor.NewPropertyValue("/tmp/mgo"))
	testProcessSession := new(apiMocks.ProcessSession)
	flowFile := new(apiMocks.FlowFile)
	testProcessSession.On("Create").Return(flowFile)
	testProcessSession.On("PutAllAttributes", flowFile, mock.AnythingOfType("map[string]string")).Return(flowFile)
	testProcessSession.On("ImportFrom", "/tmp/mgo/testFile", false, flowFile).Return(flowFile)
	testProcessSession.On("Transfer", flowFile, success)

	getFile.OnTrigger(testProcessContext, testProcessSession)

	testProcessContext.AssertExpectations(t)
	testProcessSession.AssertExpectations(t)
	mockFile.AssertExpectations(t)
}
