package chess

import (
    "testing"

	"github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/assert"
)

type MockPlayerTransitionFactory struct {
	mock.Mock
}

func (m *MockPlayerTransitionFactory) newIncrementalTransitionAsPlayerTransition(b Board, p PlayerCollection) (PlayerTransition, error) {
	args := m.Called(b, p)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(PlayerTransition), args.Error(1)
	}
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

func Test_IncrementalTransition(t *testing.T) {
	board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    playerCollection.On("getCurrent").Return("white", nil)
    playerCollection.On("getWinner").Return("", nil)
    playerCollection.On("getNext").Return(&Player{"black", true}, nil)
    playerCollection.On("setCurrent", "black").Return(nil)
    playerCollection.On("setCurrent", "white").Return(nil)
    board.On("CalculateMoves", "black").Return(nil)
    board.On("Checkmate").Return(false)
    board.On("Stalemate").Return(false)
    incrementalTransition, err := playerTransitionFactoryInstance.newIncrementalTransition(board, playerCollection)
    assert.Nil(t, err)

    playerCollection.On("setCurrent", "black").Return(nil)
    playerCollection.On("setWinner", "").Return(nil)
    board.On("CalculateMoves", "black").Return(nil)
    err = incrementalTransition.execute()
    assert.Nil(t, err)

    playerCollection.On("setCurrent", "white").Return(nil)
    playerCollection.On("setWinner", "").Return(nil)
    board.On("CalculateMoves", "white").Return(nil)
    err = incrementalTransition.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
	playerCollection.AssertExpectations(t)
}

