package chess

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_IncrementalTransition_inCheckmate(t *testing.T) {
    white := 0
    black := 1

	board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    playerCollection.On("getCurrent").Return(white, true).Once()
    playerCollection.On("getWinner").Return(-1, true).Once()
    playerCollection.On("getGameOver").Return(false, nil).Once()
    playerCollection.On("getNextAndRemaining").Return(black, 2, nil).Once()
    incrementalTransition, err := createPlayerTransition(board, playerCollection, true, false)
    assert.Nil(t, err)

    playerCollection.On("setCurrent", black).Return(true)
    playerCollection.On("setWinner", black).Return(true)
    playerCollection.On("setGameOver", true).Return(nil)
    playerCollection.On("eliminate", white).Return(nil)
    board.On("disablePieces", white, true).Return(nil)
    err = incrementalTransition.execute()
    assert.Nil(t, err)

    playerCollection.On("setCurrent", white).Return(true)
    playerCollection.On("setWinner", -1).Return(true)
    playerCollection.On("setGameOver", false).Return(nil)
    playerCollection.On("restore", white).Return(nil)
    board.On("disablePieces", white, false).Return(nil)
    err = incrementalTransition.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
	playerCollection.AssertExpectations(t)
}

func Test_IncrementalTransition_noCheckmate(t *testing.T) {
    white := 0
    black := 1

	board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    playerCollection.On("getCurrent").Return(white, true)
    playerCollection.On("getWinner").Return(-1, true)
    playerCollection.On("getGameOver").Return(false, nil)
    playerCollection.On("getNextAndRemaining").Return(black, 2, nil)
    incrementalTransition, err := createPlayerTransition(board, playerCollection, false, false)
    assert.Nil(t, err)

    playerCollection.On("setCurrent", black).Return(true)
    playerCollection.On("setWinner", -1).Return(true)
    playerCollection.On("setGameOver", false).Return(nil)
    err = incrementalTransition.execute()
    assert.Nil(t, err)

    playerCollection.On("setCurrent", white).Return(true)
    playerCollection.On("setWinner", -1).Return(true)
    playerCollection.On("setGameOver", false).Return(nil)
    err = incrementalTransition.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
	playerCollection.AssertExpectations(t)
}

