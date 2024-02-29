package chess

import (
    "testing"
    "strings"
    "fmt"
    "time"

    "github.com/stretchr/testify/assert"
)

func Test_Minimax(t *testing.T) {
    b, err := newSimpleBoard(Point{4, 4})
    assert.Nil(t, err)
    b.setPiece(Point{2, 3}, newKing("white", false, 0, -1))
    b.setPiece(Point{3, 3}, newRook("white", false))
    b.setPiece(Point{0, 0}, newKing("black", false, 0, 1))
    b.disableLocation(Point{0, 3})
    err = b.CalculateMoves()
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection([]Player{{"white", true}, {"black", true}})
    assert.Nil(t, err)

    i, err := invokerFactoryInstance.newSimpleInvoker()
    assert.Nil(t, err)

    game := &SimpleGame{
        b: b,
        p: p,
        i: i,
    }

    searcher, err := newSimpleSearcher(game)
    assert.Nil(t, err)

    score, move, err := searcher.minimax(4)
    assert.Nil(t, err)

    assert.Equal(t, 100000, score["white"])
    assert.Equal(t, -100000, score["black"])
    assert.Equal(t, 1, move.getAction().toLocation.x)
    assert.Equal(t, 2, move.getAction().toLocation.y)

    actualPrintedBoard := game.Print()
    expectedPrintedBoard := strings.Trim(`
+---------------------------------------------------+
|         0x |         1x |         2x |         3x |
| K black    |            |            |            |
|         0y |         0y |         0y |         0y |
+---------------------------------------------------+
|         0x |         1x |         2x |         3x |
|            |            |            |            |
|         1y |         1y |         1y |         1y |
+---------------------------------------------------+
|         0x |         1x |         2x |         3x |
|            |            |            |            |
|         2y |         2y |         2y |         2y |
+---------------------------------------------------+
|         0x |         1x |         2x |         3x |
|XXXXXXXXXXXX|            | K white    | R white    |
|         3y |         3y |         3y |         3y |
+---------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

func Test_Minimax_AvoidMateInOne(t *testing.T) {
    b, err := createSimpleBoardWithDefaultPieceLocations()
    assert.Nil(t, err)
    b.setPiece(Point{5, 4}, newQueen("white"))
    b.setPiece(Point{5, 3}, newQueen("white"))
    b.setPiece(Point{3, 2}, newPawn("black", true, 0, 1))
    b.setPiece(Point{6, 0}, nil)
    err = b.CalculateMoves()
    assert.Nil(t, err)

    p, err := createSimplePlayerCollectionWithDefaultPlayers()
    assert.Nil(t, err)
    p.setCurrent("black")

    i, err := invokerFactoryInstance.newSimpleInvoker()
    assert.Nil(t, err)

    game := &SimpleGame{
        b: b,
        p: p,
        i: i,
    }

    searcher, err := newSimpleSearcher(game)
    assert.Nil(t, err)

    _, move, err := searcher.minimax(3)
    assert.Nil(t, err)

    assert.Equal(t, 5, move.getAction().toLocation.x)
    assert.Equal(t, 2, move.getAction().toLocation.y)

    actualPrintedBoard := game.Print()
    expectedPrintedBoard := strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R black    | N black    | B black    | Q black    | K black    | B black    |            | R black    |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P black    | P black    | P black    | P black    | P black    | P black    | P black    | P black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            | P black    |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            | Q white    |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            | Q white    |            |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P white    | P white    | P white    | P white    | P white    | P white    | P white    | P white    |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R white    | N white    | B white    | Q white    | K white    | B white    | N white    | R white    |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

// go test ./chess -bench='Benchmark_Minimax' -cpuprofile='cpu.prof' -memprofile='mem.prof' -trace='trace.out' -run NONE
// go tool pprof cpu.prof
// go tool trace trace.out
// web
// go tool pprof -http=:8000 cpu.prof
func Benchmark_Minimax(t *testing.B) {
    for i := 0; i < t.N; i++ {
        b, err := newSimpleBoard(Point{8, 8})
        assert.Nil(t, err)

        b.setPiece(Point{0, 0}, newQueen("black"))
        b.setPiece(Point{3, 0}, newKing("black", true, 0, 1))
        b.setPiece(Point{6, 0}, newQueen("black"))

        b.setPiece(Point{2, 1}, newKnight("black"))
        b.setPiece(Point{3, 1}, newQueen("black"))
        b.setPiece(Point{4, 1}, newKnight("black"))
        b.setPiece(Point{7, 1}, newBishop("black"))

        b.setPiece(Point{1, 2}, newKnight("black"))
        b.setPiece(Point{3, 2}, newPawn("white", true, 0, -1))
        b.setPiece(Point{5, 2}, newKnight("black"))
        b.setPiece(Point{7, 2}, newBishop("black"))

        b.setPiece(Point{2, 3}, newRook("black", true))
        b.setPiece(Point{3, 3}, newKnight("black"))
        b.setPiece(Point{4, 3}, newRook("black", true))
        b.setPiece(Point{7, 3}, newQueen("white"))

        b.setPiece(Point{1, 4}, newKnight("white"))
        b.setPiece(Point{2, 4}, newQueen("white"))
        b.setPiece(Point{4, 4}, newQueen("white"))
        b.setPiece(Point{5, 4}, newKnight("white"))
        b.setPiece(Point{7, 4}, newQueen("white"))

        b.setPiece(Point{3, 5}, newQueen("white"))
        b.setPiece(Point{7, 5}, newBishop("white"))

        b.setPiece(Point{2, 6}, newRook("white", true))
        b.setPiece(Point{3, 6}, newQueen("white"))
        b.setPiece(Point{4, 6}, newRook("white", true))
        b.setPiece(Point{7, 6}, newBishop("white"))

        b.setPiece(Point{0, 7}, newQueen("white"))
        b.setPiece(Point{3, 7}, newKing("white", true, 0, -1))
        b.setPiece(Point{6, 7}, newQueen("white"))

        err = b.CalculateMoves()
        assert.Nil(t, err)

        // without calculateMoves, the time falls to ~2s for depth 7
        // without calculateToLocations, the time stays the same
        // b.test = true

        p, err := newSimplePlayerCollection([]Player{{"white", true}, {"black", true}})
        assert.Nil(t, err)

        i, err := invokerFactoryInstance.newSimpleInvoker()
        assert.Nil(t, err)

        game := &SimpleGame{
            b: b,
            p: p,
            i: i,
        }

        searcher, err := newSimpleSearcher(game)
        assert.Nil(t, err)

        _, _, err = searcher.minimax(3)
        assert.Nil(t, err)

        actualPrintedBoard := game.Print()
        expectedPrintedBoard := strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| Q black    |            |            | K black    |            |            | Q black    |            |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | N black    | Q black    | N black    |            |            | B black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            | N black    |            | P white    |            | N black    |            | B black    |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | R black    | N black    | R black    |            |            | Q white    |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            | N white    | Q white    |            | Q white    | N white    |            | Q white    |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            | Q white    |            |            |            | B white    |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | R white    | Q white    | R white    |            |            | B white    |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| Q white    |            |            | K white    |            |            | Q white    |            |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
        `, " \t\n") + "\n"
        assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)

        fmt.Println("minimaxCalls", searcher.minimaxCalls)
    }
}

// without calculateMoves, the time falls to ~2s for depth 7
// without calculateToLocations, the time stays the same
// first we should try to calculate move for affected squared only then we can try more extreme solutions
func Benchmark_CalculateMoves(t *testing.B) {
    for i := 0; i < t.N; i++ {
        b, err := newSimpleBoard(Point{8, 8})
        assert.Nil(t, err)

        b.setPiece(Point{0, 0}, newQueen("black"))
        b.setPiece(Point{3, 0}, newKing("black", true, 0, 1))
        b.setPiece(Point{6, 0}, newQueen("black"))

        b.setPiece(Point{2, 1}, newKnight("black"))
        b.setPiece(Point{3, 1}, newQueen("black"))
        b.setPiece(Point{4, 1}, newKnight("black"))
        b.setPiece(Point{7, 1}, newBishop("black"))

        b.setPiece(Point{1, 2}, newKnight("black"))
        b.setPiece(Point{3, 2}, newPawn("white", true, 0, -1))
        b.setPiece(Point{5, 2}, newKnight("black"))
        b.setPiece(Point{7, 2}, newBishop("black"))

        b.setPiece(Point{2, 3}, newRook("black", true))
        b.setPiece(Point{3, 3}, newKnight("black"))
        b.setPiece(Point{4, 3}, newRook("black", true))
        b.setPiece(Point{7, 3}, newQueen("white"))

        b.setPiece(Point{1, 4}, newKnight("white"))
        b.setPiece(Point{2, 4}, newQueen("white"))
        b.setPiece(Point{4, 4}, newQueen("white"))
        b.setPiece(Point{5, 4}, newKnight("white"))
        b.setPiece(Point{7, 4}, newQueen("white"))

        b.setPiece(Point{3, 5}, newQueen("white"))
        b.setPiece(Point{7, 5}, newBishop("white"))

        b.setPiece(Point{2, 6}, newRook("white", true))
        b.setPiece(Point{3, 6}, newQueen("white"))
        b.setPiece(Point{4, 6}, newRook("white", true))
        b.setPiece(Point{7, 6}, newBishop("white"))

        b.setPiece(Point{0, 7}, newQueen("white"))
        b.setPiece(Point{3, 7}, newKing("white", true, 0, -1))
        b.setPiece(Point{6, 7}, newQueen("white"))

        // the time is around 1ms right now
        // when we skip the calculation it's around 40ns
        // when we skip the calculating checks, it's around 100us
        // b.test = true

        start := time.Now()
        err = b.CalculateMoves()
        end := time.Now()

        assert.Nil(t, err)
        fmt.Println("time", end.Sub(start))
    }
}

