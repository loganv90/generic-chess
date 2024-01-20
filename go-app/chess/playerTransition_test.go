package chess

import (
	"github.com/stretchr/testify/mock"
)

type MockPlayerTransitionFactory struct {
	mock.Mock
}

func (m *MockPlayerTransitionFactory) newIncrementalTransition(b Board, p PlayerCollection) (*IncrementalTransition, error) {
	args := m.Called(b, p)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*IncrementalTransition), args.Error(1)
	}
}

type MockPlayerTransition struct {
	mock.Mock
}

func (m *MockPlayerTransition) execute() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPlayerTransition) undo() error {
	args := m.Called()
	return args.Error(0)
}

