package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_SimpleInvoker_UndoAndRedoInOrder(t *testing.T) {
    b := &MockBoard{}
    b.On("getEnPassant", mock.Anything).Return(&EnPassant{})
    b.On("getVulnerable", mock.Anything).Return(&Vulnerable{})
    move1 := FastMove{b: b}
    move2 := FastMove{b: b}

    p := &MockPlayerCollection{}
    p.On("setCurrent", mock.Anything, mock.Anything).Return(true)
    p.On("setWinner", mock.Anything, mock.Anything).Return(true)
    p.On("setGameOver", mock.Anything, mock.Anything).Return(nil)
    playerTransition1 := PlayerTransition{b: b, p: p}
    playerTransition2 := PlayerTransition{b: b, p: p}

	simpleInvoker, err := invokerFactoryInstance.newSimpleInvoker()
	assert.Nil(t, err)

	err = simpleInvoker.execute(move1, playerTransition1)
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "getEnPassant", 1)

	err = simpleInvoker.execute(move2, playerTransition2)
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "getEnPassant", 2)

	err = simpleInvoker.undo()
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "getEnPassant", 3)

	err = simpleInvoker.undo()
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "getEnPassant", 4)

	err = simpleInvoker.redo()
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "getEnPassant", 5)

	err = simpleInvoker.redo()
	assert.Nil(t, err)
	b.AssertNumberOfCalls(t, "getEnPassant", 6)
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
    b := &MockBoard{}
    b.On("getEnPassant", mock.Anything).Return(&EnPassant{})
    b.On("getVulnerable", mock.Anything).Return(&Vulnerable{})
    move1 := FastMove{b: b}
    move2 := FastMove{b: b}
    move3 := FastMove{b: b}

    p := &MockPlayerCollection{}
    p.On("setCurrent", mock.Anything, mock.Anything).Return(true)
    p.On("setWinner", mock.Anything, mock.Anything).Return(true)
    p.On("setGameOver", mock.Anything, mock.Anything).Return(nil)
    playerTransition1 := PlayerTransition{b: b, p: p}
    playerTransition2 := PlayerTransition{b: b, p: p}
    playerTransition3 := PlayerTransition{b: b, p: p}

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

