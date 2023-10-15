package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SimpleInvoker_UndoAndRedoInOrder(t *testing.T) {
	move1 := &mockMove{}
	move1.On("execute").Return(nil)
	move1.On("undo").Return(nil)

	move2 := &mockMove{}
	move2.On("execute").Return(nil)
	move2.On("undo").Return(nil)

	simpleInvoker, err := invokerFactoryInstance.newSimpleInvoker()
	assert.Nil(t, err)

	err = simpleInvoker.execute(move1)
	assert.Nil(t, err)
	move1.AssertNumberOfCalls(t, "execute", 1)

	err = simpleInvoker.execute(move2)
	assert.Nil(t, err)
	move2.AssertNumberOfCalls(t, "execute", 1)

	err = simpleInvoker.undo()
	assert.Nil(t, err)
	move2.AssertNumberOfCalls(t, "undo", 1)

	err = simpleInvoker.undo()
	assert.Nil(t, err)
	move1.AssertNumberOfCalls(t, "undo", 1)

	err = simpleInvoker.redo()
	assert.Nil(t, err)
	move1.AssertNumberOfCalls(t, "execute", 2)

	err = simpleInvoker.redo()
	assert.Nil(t, err)
	move2.AssertNumberOfCalls(t, "execute", 2)
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
	move1 := &mockMove{}
	move1.On("execute").Return(nil)
	move1.On("undo").Return(nil)

	move2 := &mockMove{}
	move2.On("execute").Return(nil)
	move2.On("undo").Return(nil)

	move3 := &mockMove{}
	move3.On("execute").Return(nil)
	move3.On("undo").Return(nil)

	simpleInvoker, err := invokerFactoryInstance.newSimpleInvoker()
	assert.Nil(t, err)

	err = simpleInvoker.execute(move1)
	assert.Nil(t, err)
	err = simpleInvoker.execute(move2)
	assert.Nil(t, err)

	err = simpleInvoker.undo()
	assert.Nil(t, err)
	err = simpleInvoker.undo()
	assert.Nil(t, err)

	err = simpleInvoker.execute(move3)
	assert.Nil(t, err)
	err = simpleInvoker.redo()
	assert.NotNil(t, err)
}