package api

import (
	"github.com/stretchr/testify/mock"
)

type MockFlowFile struct {
	mock.Mock
}

func (m *MockFlowFile) Id() uint64 {
	args := m.Called()
	return args.Get(0).(uint64)
}
