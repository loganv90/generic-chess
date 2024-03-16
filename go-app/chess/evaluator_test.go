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

            pieceLocations := make([]Array100[Point], 4)


            pieceLocations[white] = Array100[Point]{}

            p := pieceLocations[white].get()
            p.x = 0
            p.y = 0
            pieceLocations[white].next()

            p = pieceLocations[white].get()
            p.x = 0
            p.y = 1
            pieceLocations[white].next()

            p = pieceLocations[white].get()
            p.x = 0
            p.y = 2
            pieceLocations[white].next()


            pieceLocations[black] = Array100[Point]{}

            p = pieceLocations[black].get()
            p.x = 1
            p.y = 0
            pieceLocations[black].next()

            p = pieceLocations[black].get()
            p.x = 1
            p.y = 1
            pieceLocations[black].next()

            p = pieceLocations[black].get()
            p.x = 1
            p.y = 2
            pieceLocations[black].next()


            pieceLocations[red] = Array100[Point]{}

            p = pieceLocations[red].get()
            p.x = 2
            p.y = 0
            pieceLocations[red].next()

            p = pieceLocations[red].get()
            p.x = 2
            p.y = 1
            pieceLocations[red].next()

            p = pieceLocations[red].get()
            p.x = 2
            p.y = 2
            pieceLocations[red].next()


            pieceLocations[blue] = Array100[Point]{}

            p = pieceLocations[blue].get()
            p.x = 3
            p.y = 0
            pieceLocations[blue].next()

            p = pieceLocations[blue].get()
            p.x = 3
            p.y = 1
            pieceLocations[blue].next()

            p = pieceLocations[blue].get()
            p.x = 3
            p.y = 2
            pieceLocations[blue].next()


            whitePiece := Piece{white, test.whiteIndex}
            board.On("getPiece", &Point{0, 0}).Return(&whitePiece)
            board.On("getPiece", &Point{0, 1}).Return(&whitePiece)
            board.On("getPiece", &Point{0, 2}).Return(&whitePiece)

            blackPiece := Piece{black, ROOK}
            board.On("getPiece", &Point{1, 0}).Return(&blackPiece)
            board.On("getPiece", &Point{1, 1}).Return(&blackPiece)
            board.On("getPiece", &Point{1, 2}).Return(&blackPiece)

            redPiece := Piece{red, KNIGHT}
            board.On("getPiece", &Point{2, 0}).Return(&redPiece)
            board.On("getPiece", &Point{2, 1}).Return(&redPiece)
            board.On("getPiece", &Point{2, 2}).Return(&redPiece)

            bluePiece := Piece{blue, PAWN_D}
            board.On("getPiece", &Point{3, 0}).Return(&bluePiece)
            board.On("getPiece", &Point{3, 1}).Return(&bluePiece)
            board.On("getPiece", &Point{3, 2}).Return(&bluePiece)

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

