package chess

import (
    "testing"
    "fmt"

    "github.com/stretchr/testify/assert"
)

func Test_Eval_Draw(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection(2)
    assert.Nil(t, err)
    p.setGameOver(true)
    p.setWinner(-1)

    evaluator := newSimpleEvaluator(b, p)

    score := make([]int, 2)
    evaluator.eval(score)
    assert.Equal(t, 0, score[white])
    assert.Equal(t, 0, score[black])
}

func Test_Eval_Win(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection(2)
    assert.Nil(t, err)
    p.setGameOver(true)
    p.setWinner(white)

    evaluator := newSimpleEvaluator(b, p)

    score := make([]int, 2)
    evaluator.eval(score)
    assert.Equal(t, 100000, score[white])
    assert.Equal(t, -100000, score[black])
}

func Test_Eval_Lose(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection(2)
    assert.Nil(t, err)
    p.setGameOver(true)
    p.setWinner(black)

    evaluator := newSimpleEvaluator(b, p)

    score := make([]int, 2)
    evaluator.eval(score)
    assert.Equal(t, -100000, score[white])
    assert.Equal(t, 100000, score[black])
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
        {QUEEN, 500, 270, 160, 50},
        {ROOK, 350, 350, 210, 70},
        {KNIGHT, 250, 410, 250, 80},
    }

    for _, test := range tests {
        testname := fmt.Sprintf("%d", test.whiteIndex)
        t.Run(testname, func(t *testing.T) {
            b, err := newSimpleBoard(10, 10, 4)
            assert.Nil(t, err)
            b.setPiece(b.getIndex(0, 0), b.getAllPiece(white, test.whiteIndex))
            b.setPiece(b.getIndex(0, 1), b.getAllPiece(white, test.whiteIndex))
            b.setPiece(b.getIndex(0, 2), b.getAllPiece(white, test.whiteIndex))
            b.setPiece(b.getIndex(1, 0), b.getAllPiece(black, ROOK))
            b.setPiece(b.getIndex(1, 1), b.getAllPiece(black, ROOK))
            b.setPiece(b.getIndex(1, 2), b.getAllPiece(black, ROOK))
            b.setPiece(b.getIndex(2, 0), b.getAllPiece(red, KNIGHT))
            b.setPiece(b.getIndex(2, 1), b.getAllPiece(red, KNIGHT))
            b.setPiece(b.getIndex(2, 2), b.getAllPiece(red, KNIGHT))
            b.setPiece(b.getIndex(3, 0), b.getAllPiece(blue, PAWN_D))
            b.setPiece(b.getIndex(3, 1), b.getAllPiece(blue, PAWN_D))
            b.setPiece(b.getIndex(3, 2), b.getAllPiece(blue, PAWN_D))

            p, err := newSimplePlayerCollection(4)
            assert.Nil(t, err)

            evaluator := newSimpleEvaluator(b, p)

            score := make([]int, 4)
            evaluator.evalMaterial(score)
            assert.Equal(t, test.whiteScore, score[white])
            assert.Equal(t, test.blackScore, score[black])
            assert.Equal(t, test.redScore, score[red])
            assert.Equal(t, test.blueScore, score[blue])
        })
    }
}

func Test_EvalMobility(t *testing.T) {
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
        {QUEEN, 0, 75, 13, 11},
        {ROOK, 43, 43, 7, 6},
        {KNIGHT, 11, 66, 11, 9},
    }

    for _, test := range tests {
        testname := fmt.Sprintf("%d", test.whiteIndex)
        t.Run(testname, func(t *testing.T) {
            b, err := newSimpleBoard(10, 10, 4)
            assert.Nil(t, err)
            b.setPiece(b.getIndex(0, 0), b.getAllPiece(white, test.whiteIndex))
            b.setPiece(b.getIndex(1, 1), b.getAllPiece(white, test.whiteIndex))
            b.setPiece(b.getIndex(0, 9), b.getAllPiece(black, ROOK))
            b.setPiece(b.getIndex(1, 8), b.getAllPiece(black, ROOK))
            b.setPiece(b.getIndex(9, 9), b.getAllPiece(red, KNIGHT))
            b.setPiece(b.getIndex(8, 8), b.getAllPiece(red, KNIGHT))
            b.setPiece(b.getIndex(9, 0), b.getAllPiece(blue, PAWN_D))
            b.setPiece(b.getIndex(8, 1), b.getAllPiece(blue, PAWN_D))
            b.CalculateMoves()

            p, err := newSimplePlayerCollection(4)
            assert.Nil(t, err)

            evaluator := newSimpleEvaluator(b, p)

            score := make([]int, 4)
            evaluator.evalMobility(score)
            assert.Equal(t, test.whiteScore, score[white])
            assert.Equal(t, test.blackScore, score[black])
            assert.Equal(t, test.redScore, score[red])
            assert.Equal(t, test.blueScore, score[blue])
        })
    }
}

