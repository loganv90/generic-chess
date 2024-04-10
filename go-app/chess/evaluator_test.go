package chess

import (
    "testing"
    "fmt"
    "math"

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
    assert.Equal(t, math.MaxInt, score[white])
    assert.Equal(t, math.MinInt, score[black])
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
    assert.Equal(t, math.MinInt, score[white])
    assert.Equal(t, math.MaxInt, score[black])
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
        {QUEEN, 2700, 1500, 900, 300},
        {ROOK, 1500, 1500, 900, 300},
        {KNIGHT, 900, 1500, 900, 300},
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

            evaluator.evalMaterial()
            assert.Equal(t, test.whiteScore, evaluator.material[white])
            assert.Equal(t, test.blackScore, evaluator.material[black])
            assert.Equal(t, test.redScore, evaluator.material[red])
            assert.Equal(t, test.blueScore, evaluator.material[blue])
        })
    }
}

func Test_EvalPosition(t *testing.T) {
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
        {QUEEN, 0, 0, 0, 100},
        {PAWN_U, 100, 0, 0, 100},
        {PAWN_D, 0, 0, 0, 100},
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
            b.setPiece(b.getIndex(9, 9), b.getAllPiece(red, PAWN_L))
            b.setPiece(b.getIndex(8, 8), b.getAllPiece(red, PAWN_L))
            b.setPiece(b.getIndex(9, 0), b.getAllPiece(blue, PAWN_R))
            b.setPiece(b.getIndex(8, 1), b.getAllPiece(blue, PAWN_R))
            b.populatePieceSquareTables()

            p, err := newSimplePlayerCollection(4)
            assert.Nil(t, err)

            evaluator := newSimpleEvaluator(b, p)

            evaluator.evalPosition()
            assert.Equal(t, test.whiteScore, evaluator.position[white])
            assert.Equal(t, test.blackScore, evaluator.position[black])
            assert.Equal(t, test.redScore, evaluator.position[red])
            assert.Equal(t, test.blueScore, evaluator.position[blue])
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
        {QUEEN, 0, 34, 6, 5},
        {ROOK, 34, 34, 6, 5},
        {KNIGHT, 6, 34, 6, 5},
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

            evaluator.evalMobility()
            assert.Equal(t, test.whiteScore, evaluator.mobility[white])
            assert.Equal(t, test.blackScore, evaluator.mobility[black])
            assert.Equal(t, test.redScore, evaluator.mobility[red])
            assert.Equal(t, test.blueScore, evaluator.mobility[blue])
        })
    }
}

