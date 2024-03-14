package chess

import (
    "testing"
    "fmt"

    "github.com/stretchr/testify/assert"
)

func Test_Eval_Draw(t *testing.T) {
    white := 0
    black := 1

    board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    evaluator, err := newSimpleEvaluator(board, playerCollection)
    assert.Nil(t, err)

    playerCollection.On("getGameOver").Return(true, nil)
    playerCollection.On("getWinner").Return(-1, true)
    playerCollection.On("getPlayers").Return(2)

    score, err := evaluator.eval()
    assert.Nil(t, err)
    assert.Equal(t, 0, score[white])
    assert.Equal(t, 0, score[black])

    playerCollection.AssertExpectations(t)
}

func Test_Eval_Win(t *testing.T) {
    white := 0
    black := 1

    board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    evaluator, err := newSimpleEvaluator(board, playerCollection)
    assert.Nil(t, err)

    playerCollection.On("getGameOver").Return(true, nil)
    playerCollection.On("getWinner").Return(white, true)
    playerCollection.On("getPlayers").Return(2, nil)

    score, err := evaluator.eval()
    assert.Nil(t, err)
    assert.Equal(t, 100000, score[white])
    assert.Equal(t, -100000, score[black])

    playerCollection.AssertExpectations(t)
}

func Test_Eval_Lose(t *testing.T) {
    white := 0
    black := 1

    board := &MockBoard{}
    playerCollection := &MockPlayerCollection{}

    evaluator, err := newSimpleEvaluator(board, playerCollection)
    assert.Nil(t, err)

    playerCollection.On("getGameOver").Return(true, nil)
    playerCollection.On("getWinner").Return(black, true)
    playerCollection.On("getPlayers").Return(2, nil)

    score, err := evaluator.eval()
    assert.Nil(t, err)
    assert.Equal(t, -100000, score[white])
    assert.Equal(t, 100000, score[black])

    playerCollection.AssertExpectations(t)
}

func Test_EvalMaterial(t *testing.T) {
    white := 0
    black := 1
    red := 2
    blue := 3

    tests := []struct {
        whiteIndex int
        whiteScore int
        blackScore int
        redScore int
        blueScore int
    }{
        {QUEEN, 0, -1200, -1800, -2400},
        {ROOK, 0, 0, -600, -1200},
        {KNIGHT, -600, 0, -600, -1200},
    }

    for _, test := range tests {
        testname := fmt.Sprintf("%d", test.whiteIndex)
        t.Run(testname, func(t *testing.T) {
            board := &MockBoard{}
            playerCollection := &MockPlayerCollection{}

            evaluator, err := newSimpleEvaluator(board, playerCollection)
            assert.Nil(t, err)

            playerCollection.On("getPlayers").Return(4)
            pieceLocations := [][]Point{
                {{0,0}, {0,1}, {0,2}}, // white
                {{1,0}, {1,1}, {1,2}}, // black
                {{2,0}, {2,1}, {2,2}}, // red
                {{3,0}, {3,1}, {3,2}}, // blue
            }

            whitePiece := Piece{white, test.whiteIndex}
            board.On("getPiece", Point{0, 0}).Return(whitePiece, true)
            board.On("getPiece", Point{0, 1}).Return(whitePiece, true)
            board.On("getPiece", Point{0, 2}).Return(whitePiece, true)

            blackPiece := Piece{black, ROOK}
            board.On("getPiece", Point{1, 0}).Return(blackPiece, true)
            board.On("getPiece", Point{1, 1}).Return(blackPiece, true)
            board.On("getPiece", Point{1, 2}).Return(blackPiece, true)

            redPiece := Piece{red, KNIGHT}
            board.On("getPiece", Point{2, 0}).Return(redPiece, true)
            board.On("getPiece", Point{2, 1}).Return(redPiece, true)
            board.On("getPiece", Point{2, 2}).Return(redPiece, true)

            bluePiece := Piece{blue, PAWN_D}
            board.On("getPiece", Point{3, 0}).Return(bluePiece, true)
            board.On("getPiece", Point{3, 1}).Return(bluePiece, true)
            board.On("getPiece", Point{3, 2}).Return(bluePiece, true)

            score, err := evaluator.evalMaterial(pieceLocations)
            assert.Nil(t, err)
            assert.Equal(t, test.whiteScore, score[white])
            assert.Equal(t, test.blackScore, score[black])
            assert.Equal(t, test.redScore, score[red])
            assert.Equal(t, test.blueScore, score[blue])

            board.AssertExpectations(t)
        })
    }
}

