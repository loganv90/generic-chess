package chess

import (
    "testing"

	"github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/assert"
)

type MockPlayerTransitionFactory struct {
	mock.Mock
}

func (m *MockPlayerTransitionFactory) newIncrementalTransitionAsPlayerTransition(b Board, p PlayerCollection, inCheckmate bool, inStalemate bool) (PlayerTransition, error) {
    args := m.Called(b, p, inCheckmate, inStalemate)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(PlayerTransition), args.Error(1)
	}
}

func (m *MockPlayerTransitionFactory) newIncrementalTransition(b Board, p PlayerCollection, inCheckmate bool, inStalemate bool) (*IncrementalTransition, error) {
    args := m.Called(b, p, inCheckmate, inStalemate)

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

func Test_IncrementalTransition_inCheckmate(t *testing.T) {
	board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    playerCollection.On("getCurrent").Return("white", nil).Once()
    playerCollection.On("getWinner").Return("", nil).Once()
    playerCollection.On("getGameOver").Return(false, nil).Once()
    playerCollection.On("getNext").Return([]Player{{"black", true}, {"white", true}}, nil).Once()
    incrementalTransition, err := playerTransitionFactoryInstance.newIncrementalTransition(board, playerCollection, true, false)
    assert.Nil(t, err)

    playerCollection.On("setCurrent", "black").Return(nil)
    playerCollection.On("setWinner", "black").Return(nil)
    playerCollection.On("setGameOver", true).Return(nil)
    playerCollection.On("eliminate", "white").Return(nil)
    board.On("disablePieces", "white", true).Return(nil)
    err = incrementalTransition.execute()
    assert.Nil(t, err)

    playerCollection.On("setCurrent", "white").Return(nil)
    playerCollection.On("setWinner", "").Return(nil)
    playerCollection.On("setGameOver", false).Return(nil)
    playerCollection.On("restore", "white").Return(nil)
    board.On("disablePieces", "white", false).Return(nil)
    err = incrementalTransition.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
	playerCollection.AssertExpectations(t)
}

func Test_IncrementalTransition_noCheckmate(t *testing.T) {
	board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    playerCollection.On("getCurrent").Return("white", nil)
    playerCollection.On("getWinner").Return("", nil)
    playerCollection.On("getGameOver").Return(false, nil)
    playerCollection.On("getNext").Return([]Player{{"black", true}, {"white", true}}, nil)
    incrementalTransition, err := playerTransitionFactoryInstance.newIncrementalTransition(board, playerCollection, false, false)
    assert.Nil(t, err)

    playerCollection.On("setCurrent", "black").Return(nil)
    playerCollection.On("setWinner", "").Return(nil)
    playerCollection.On("setGameOver", false).Return(nil)
    err = incrementalTransition.execute()
    assert.Nil(t, err)

    playerCollection.On("setCurrent", "white").Return(nil)
    playerCollection.On("setWinner", "").Return(nil)
    playerCollection.On("setGameOver", false).Return(nil)
    err = incrementalTransition.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
	playerCollection.AssertExpectations(t)
}

