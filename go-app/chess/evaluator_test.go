package chess

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_Eval_Draw(t *testing.T) {
    board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    evaluator, err := newSimpleEvaluator(board, playerCollection)
    assert.Nil(t, err)

    playerCollection.On("getGameOver").Return(true, nil)
    playerCollection.On("getWinner").Return("", nil)
    playerCollection.On("getCurrent").Return("white", nil)

    score, err := evaluator.eval()
    assert.Nil(t, err)
    assert.Equal(t, 0, score)

    playerCollection.AssertExpectations(t)
}

func Test_Eval_Win(t *testing.T) {
    board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    evaluator, err := newSimpleEvaluator(board, playerCollection)
    assert.Nil(t, err)

    playerCollection.On("getGameOver").Return(true, nil)
    playerCollection.On("getWinner").Return("white", nil)
    playerCollection.On("getCurrent").Return("white", nil)

    score, err := evaluator.eval()
    assert.Nil(t, err)
    assert.Equal(t, 100000, score)

    playerCollection.AssertExpectations(t)
}

func Test_Eval_Lose(t *testing.T) {
    board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    evaluator, err := newSimpleEvaluator(board, playerCollection)
    assert.Nil(t, err)

    playerCollection.On("getGameOver").Return(true, nil)
    playerCollection.On("getWinner").Return("black", nil)
    playerCollection.On("getCurrent").Return("white", nil)

    score, err := evaluator.eval()
    assert.Nil(t, err)
    assert.Equal(t, -100000, score)

    playerCollection.AssertExpectations(t)
}

