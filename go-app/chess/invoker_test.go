package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_SimpleInvoker_UndoAndRedoInOrder(t *testing.T) {
    t.Cleanup(func() { playerTransitionFactoryInstance = &ConcretePlayerTransitionFactory{} })

    b := &MockBoard{}
    b.On("setEnPassant", mock.Anything, mock.Anything).Return(nil)
    b.On("setVulnerable", mock.Anything, mock.Anything).Return(nil)
    move1 := FastMove{b: b}
    move2 := FastMove{b: b}

    playerTransition1 := &MockPlayerTransition{}
    playerTransition1.On("execute").Return(nil)
    playerTransition1.On("undo").Return(nil)

    playerTransition2 := &MockPlayerTransition{}
    playerTransition2.On("execute").Return(nil)
    playerTransition2.On("undo").Return(nil)

    board1 := &MockBoard{}
    board1.On("CalculateMoves").Return(nil)
    playerCollection1 := &MockPlayerCollection{}

    board2 := &MockBoard{}
    board2.On("CalculateMoves").Return(nil)
    playerCollection2 := &MockPlayerCollection{}

	playerTransitionFactory := &MockPlayerTransitionFactory{}
	playerTransitionFactory.On("newIncrementalTransitionAsPlayerTransition", board1, playerCollection1).Return(playerTransition1, nil)
	playerTransitionFactory.On("newIncrementalTransitionAsPlayerTransition", board2, playerCollection2).Return(playerTransition2, nil)
	playerTransitionFactoryInstance = playerTransitionFactory

	simpleInvoker, err := invokerFactoryInstance.newSimpleInvoker()
	assert.Nil(t, err)

	err = simpleInvoker.execute(move1, playerTransition1)
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "setEnPassant", 1)

	err = simpleInvoker.execute(move2, playerTransition2)
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "setEnPassant", 2)

	err = simpleInvoker.undo()
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "setEnPassant", 3)

	err = simpleInvoker.undo()
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "setEnPassant", 4)

	err = simpleInvoker.redo()
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "setEnPassant", 5)

	err = simpleInvoker.redo()
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "setEnPassant", 6)
}

func Test_SimpleInvoker_UndoAndRedoWithNoMoves(t *testing.T) {
	simpleInvoker, err := invokerFactoryInstance.newSimpleInvoker()
	assert.Nil(t, err)

	err = simpleInvoker.undo()
	assert.NotNil(t, err)

	err = simpleInvoker.redo()
	assert.NotNil(t, err)
}

func Test_SimpleInvoker_OverwriteHistory(t *testing.T) {
    t.Cleanup(func() { playerTransitionFactoryInstance = &ConcretePlayerTransitionFactory{} })

    b := &MockBoard{}
    b.On("setEnPassant", mock.Anything, mock.Anything).Return(nil)
    b.On("setVulnerable", mock.Anything, mock.Anything).Return(nil)
    move1 := FastMove{b: b}
    move2 := FastMove{b: b}
    move3 := FastMove{b: b}

    playerTransition1 := &MockPlayerTransition{}
    playerTransition1.On("execute").Return(nil)
    playerTransition1.On("undo").Return(nil)

    playerTransition2 := &MockPlayerTransition{}
    playerTransition2.On("execute").Return(nil)
    playerTransition2.On("undo").Return(nil)

    playerTransition3 := &MockPlayerTransition{}
    playerTransition3.On("execute").Return(nil)
    playerTransition3.On("undo").Return(nil)

    board1 := &MockBoard{}
    board1.On("CalculateMoves").Return(nil)
    playerCollection1 := &MockPlayerCollection{}

    board2 := &MockBoard{}
    board2.On("CalculateMoves").Return(nil)
    playerCollection2 := &MockPlayerCollection{}

    board3 := &MockBoard{}
    board3.On("CalculateMoves").Return(nil)
    playerCollection3 := &MockPlayerCollection{}

	playerTransitionFactory := &MockPlayerTransitionFactory{}
	playerTransitionFactory.On("newIncrementalTransitionAsPlayerTransition", board1, playerCollection1).Return(playerTransition1, nil)
	playerTransitionFactory.On("newIncrementalTransitionAsPlayerTransition", board2, playerCollection2).Return(playerTransition2, nil)
	playerTransitionFactory.On("newIncrementalTransitionAsPlayerTransition", board3, playerCollection3).Return(playerTransition3, nil)
	playerTransitionFactoryInstance = playerTransitionFactory

	simpleInvoker, err := invokerFactoryInstance.newSimpleInvoker()
	assert.Nil(t, err)

	err = simpleInvoker.execute(move1, playerTransition1)
	assert.Nil(t, err)
	err = simpleInvoker.execute(move2, playerTransition2)
	assert.Nil(t, err)

	err = simpleInvoker.undo()
	assert.Nil(t, err)
	err = simpleInvoker.undo()
	assert.Nil(t, err)

	err = simpleInvoker.execute(move3, playerTransition3)
	assert.Nil(t, err)
	err = simpleInvoker.redo()
	assert.NotNil(t, err)
}

