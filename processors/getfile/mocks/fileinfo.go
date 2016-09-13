package mocks

import (
	"os"
	"testing"
	"time"

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
