package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SimpleInvoker_UndoAndRedoInOrder(t *testing.T) {
    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)
    move1 := FastMove{b: b}
    move2 := FastMove{b: b}

    p, err := newSimplePlayerCollection(2)
    assert.Nil(t, err)
    playerTransition1 := PlayerTransition{b: b, p: p}
    playerTransition2 := PlayerTransition{b: b, p: p}

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

	err = simpleInvoker.redo()
	assert.Nil(t, err)
	err = simpleInvoker.redo()
	assert.Nil(t, err)
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
    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)
    move1 := FastMove{b: b}
    move2 := FastMove{b: b}
    move3 := FastMove{b: b}

    p, err := newSimplePlayerCollection(2)
    assert.Nil(t, err)
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

