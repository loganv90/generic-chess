package chess

import (
    "testing"
    "fmt"

    "github.com/stretchr/testify/assert"
)

func Test_Eval_Draw(t *testing.T) {
    board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    evaluator, err := newSimpleEvaluator(board, playerCollection)
    assert.Nil(t, err)

    playerCollection.On("getGameOver").Return(true, nil)
    playerCollection.On("getWinner").Return("", nil)

    score, err := evaluator.eval("white")
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

    score, err := evaluator.eval("white")
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

    score, err := evaluator.eval("white")
    assert.Nil(t, err)
    assert.Equal(t, -100000, score)

    playerCollection.AssertExpectations(t)
}

func Test_EvalMaterial(t *testing.T) {
    tests := []struct {
        whiteValue int
        score int
    }{
        {500, 300},
        {100, -900},
    }

    for _, test := range tests {
        testname := fmt.Sprintf("%d_%d", test.whiteValue, test.score)
        t.Run(testname, func(t *testing.T) {
            board := &MockBoard{}
            playerCollection := &MockPlayerCollection{}

            evaluator, err := newSimpleEvaluator(board, playerCollection)
            assert.Nil(t, err)

            pieceLocations := map[string][]*Point{
                "white": {{0,0}, {0,1}, {0,2}},
                "black": {{1,0}, {1,1}, {1,2}},
                "red": {{2,0}, {2,1}, {2,2}},
                "blue": {{3,0}, {3,1}, {3,2}},
            }

            whitePiece := &MockPiece{}
            whitePiece.On("getValue").Return(test.whiteValue)
            board.On("getPiece", &Point{0, 0}).Return(whitePiece, nil)
            board.On("getPiece", &Point{0, 1}).Return(whitePiece, nil)
            board.On("getPiece", &Point{0, 2}).Return(whitePiece, nil)

            blackPiece := &MockPiece{}
            blackPiece.On("getValue").Return(400)
            board.On("getPiece", &Point{1, 0}).Return(blackPiece, nil)
            board.On("getPiece", &Point{1, 1}).Return(blackPiece, nil)
            board.On("getPiece", &Point{1, 2}).Return(blackPiece, nil)

            redPiece := &MockPiece{}
            redPiece.On("getValue").Return(300)
            board.On("getPiece", &Point{2, 0}).Return(redPiece, nil)
            board.On("getPiece", &Point{2, 1}).Return(redPiece, nil)
            board.On("getPiece", &Point{2, 2}).Return(redPiece, nil)

            bluePiece := &MockPiece{}
            bluePiece.On("getValue").Return(200)
            board.On("getPiece", &Point{3, 0}).Return(bluePiece, nil)
            board.On("getPiece", &Point{3, 1}).Return(bluePiece, nil)
            board.On("getPiece", &Point{3, 2}).Return(bluePiece, nil)

            score, err := evaluator.evalMaterial("white", pieceLocations)
            assert.Nil(t, err)
            assert.Equal(t, test.score, score)

            board.AssertExpectations(t)
            whitePiece.AssertExpectations(t)
            blackPiece.AssertExpectations(t)
            redPiece.AssertExpectations(t)
            bluePiece.AssertExpectations(t)
        })
    }
}

