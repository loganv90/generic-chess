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
    playerCollection.On("getPlayerColors").Return([]string{"white", "black"}, nil)

    score, err := evaluator.eval()
    assert.Nil(t, err)
    assert.Equal(t, 0, score["white"])
    assert.Equal(t, 0, score["black"])

    playerCollection.AssertExpectations(t)
}

func Test_Eval_Win(t *testing.T) {
    board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    evaluator, err := newSimpleEvaluator(board, playerCollection)
    assert.Nil(t, err)

    playerCollection.On("getGameOver").Return(true, nil)
    playerCollection.On("getWinner").Return("white", nil)
    playerCollection.On("getPlayerColors").Return([]string{"white", "black"}, nil)

    score, err := evaluator.eval()
    assert.Nil(t, err)
    assert.Equal(t, 100000, score["white"])
    assert.Equal(t, -100000, score["black"])

    playerCollection.AssertExpectations(t)
}

func Test_Eval_Lose(t *testing.T) {
    board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    evaluator, err := newSimpleEvaluator(board, playerCollection)
    assert.Nil(t, err)

    playerCollection.On("getGameOver").Return(true, nil)
    playerCollection.On("getWinner").Return("black", nil)
    playerCollection.On("getPlayerColors").Return([]string{"white", "black"}, nil)

    score, err := evaluator.eval()
    assert.Nil(t, err)
    assert.Equal(t, -100000, score["white"])
    assert.Equal(t, 100000, score["black"])

    playerCollection.AssertExpectations(t)
}

func Test_EvalMaterial(t *testing.T) {
    tests := []struct {
        whiteValue int
        whiteScore int
        blackScore int
        redScore int
        blueScore int
    }{
        {500, 0, -300, -600, -900},
        {100, -900, 0, -300, -600},
    }

    for _, test := range tests {
        testname := fmt.Sprintf("%d", test.whiteValue)
        t.Run(testname, func(t *testing.T) {
            board := &MockBoard{}
            playerCollection := &MockPlayerCollection{}

            evaluator, err := newSimpleEvaluator(board, playerCollection)
            assert.Nil(t, err)

            playerCollection.On("getPlayerColors").Return([]string{"white", "black", "red", "blue"}, nil)
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

            score, err := evaluator.evalMaterial(pieceLocations)
            assert.Nil(t, err)
            assert.Equal(t, test.whiteScore, score["white"])
            assert.Equal(t, test.blackScore, score["black"])
            assert.Equal(t, test.redScore, score["red"])
            assert.Equal(t, test.blueScore, score["blue"])

            board.AssertExpectations(t)
            whitePiece.AssertExpectations(t)
            blackPiece.AssertExpectations(t)
            redPiece.AssertExpectations(t)
            bluePiece.AssertExpectations(t)
        })
    }
}

