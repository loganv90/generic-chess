package chess

import (
    "testing"
    "strings"
    "fmt"
    "time"

    "github.com/stretchr/testify/assert"
)

func Test_Minimax(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(Point{4, 4}, 2)
    assert.Nil(t, err)
    setPiece(b, Point{2, 3}, Piece{white, KING_U})
    setPiece(b, Point{3, 3}, Piece{white, ROOK})
    setPiece(b, Point{0, 0}, Piece{black, KING_D})
    b.disableLocation(&Point{0, 3})
    b.CalculateMoves()

    p, err := newSimplePlayerCollection(2)
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

    score, move, _, err := searcher.minimax(4)
    assert.Nil(t, err)

    assert.Equal(t, 100000, score[white])
    assert.Equal(t, -100000, score[black])
    assert.Equal(t, 1, move.toLocation.x)
    assert.Equal(t, 2, move.toLocation.y)

    actualPrintedBoard := game.Print()
    expectedPrintedBoard := strings.Trim(`
+---------------------------------------------------+
|         0x |         1x |         2x |         3x |
| K 1        |            |            |            |
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
|XXXXXXXXXXXX|            | K 0        | R 0        |
|         3y |         3y |         3y |         3y |
+---------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

func Test_Minimax_AvoidMateInOne(t *testing.T) {
    white := 0
    black := 1

    b, err := createSimpleBoardWithDefaultPieceLocations()
    assert.Nil(t, err)
    setPiece(b, Point{5, 4}, Piece{white, QUEEN})
    setPiece(b, Point{5, 3}, Piece{white, QUEEN})
    setPiece(b, Point{3, 2}, Piece{black, PAWN_D_M})
    setPiece(b, Point{6, 0}, Piece{0, 0})
    b.CalculateMoves()

    p, err := createSimplePlayerCollectionWithDefaultPlayers()
    assert.Nil(t, err)
    p.setCurrent(black)

    i, err := invokerFactoryInstance.newSimpleInvoker()
    assert.Nil(t, err)

    game := &SimpleGame{
        b: b,
        p: p,
        i: i,
    }

    searcher, err := newSimpleSearcher(game)
    assert.Nil(t, err)

    what, move, ok, err := searcher.minimax(3)
    assert.Nil(t, err)

    fmt.Println("ok", ok)
    fmt.Println("err", err)
    fmt.Println("move", move)
    fmt.Println("what", what)

    assert.Equal(t, 5, move.toLocation.x)
    assert.Equal(t, 2, move.toLocation.y)

    actualPrintedBoard := game.Print()
    expectedPrintedBoard := strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R 1        | N 1        | B 1        | Q 1        | K 1        | B 1        |            | R 1        |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P 1        | P 1        | P 1        | P 1        | P 1        | P 1        | P 1        | P 1        |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            | P 1        |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            | Q 0        |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            | Q 0        |            |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P 0        | P 0        | P 0        | P 0        | P 0        | P 0        | P 0        | P 0        |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R 0        | N 0        | B 0        | Q 0        | K 0        | B 0        | N 0        | R 0        |
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
    white := 0
    black := 1

    for i := 0; i < t.N; i++ {
        b, err := newSimpleBoard(Point{8, 8}, 2)
        assert.Nil(t, err)

        setPiece(b, Point{0, 0}, Piece{black, QUEEN})
        setPiece(b, Point{3, 0}, Piece{black, KING_U})
        setPiece(b, Point{6, 0}, Piece{black, QUEEN})

        setPiece(b, Point{2, 1}, Piece{black, KNIGHT})
        setPiece(b, Point{3, 1}, Piece{black, QUEEN})
        setPiece(b, Point{4, 1}, Piece{black, KNIGHT})
        setPiece(b, Point{7, 1}, Piece{black, BISHOP})

        setPiece(b, Point{1, 2}, Piece{black, KNIGHT})
        setPiece(b, Point{3, 2}, Piece{white, PAWN_U})
        setPiece(b, Point{5, 2}, Piece{black, KNIGHT})
        setPiece(b, Point{7, 2}, Piece{black, BISHOP})

        setPiece(b, Point{2, 3}, Piece{black, ROOK})
        setPiece(b, Point{3, 3}, Piece{black, KNIGHT})
        setPiece(b, Point{4, 3}, Piece{black, ROOK})
        setPiece(b, Point{7, 3}, Piece{white, QUEEN})

        setPiece(b, Point{1, 4}, Piece{white, KNIGHT})
        setPiece(b, Point{2, 4}, Piece{white, QUEEN})
        setPiece(b, Point{4, 4}, Piece{white, QUEEN})
        setPiece(b, Point{5, 4}, Piece{white, KNIGHT})
        setPiece(b, Point{7, 4}, Piece{white, QUEEN})

        setPiece(b, Point{3, 5}, Piece{white, QUEEN})
        setPiece(b, Point{7, 5}, Piece{white, BISHOP})

        setPiece(b, Point{2, 6}, Piece{white, ROOK})
        setPiece(b, Point{3, 6}, Piece{white, QUEEN})
        setPiece(b, Point{4, 6}, Piece{white, ROOK})
        setPiece(b, Point{7, 6}, Piece{white, BISHOP})

        setPiece(b, Point{0, 7}, Piece{white, QUEEN})
        setPiece(b, Point{3, 7}, Piece{white, KING_U})
        setPiece(b, Point{6, 7}, Piece{white, QUEEN})

        b.CalculateMoves()

        // without calculateMoves, the time falls to ~2s for depth 7
        // without calculateToLocations, the time stays the same
        // b.test = true

        p, err := newSimplePlayerCollection(2)
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

        _, _, _, err = searcher.minimax(4)
        assert.Nil(t, err)

        actualPrintedBoard := game.Print()
        expectedPrintedBoard := strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| Q 1        |            |            | K 1        |            |            | Q 1        |            |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | N 1        | Q 1        | N 1        |            |            | B 1        |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            | N 1        |            | P 0        |            | N 1        |            | B 1        |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | R 1        | N 1        | R 1        |            |            | Q 0        |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            | N 0        | Q 0        |            | Q 0        | N 0        |            | Q 0        |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            | Q 0        |            |            |            | B 0        |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | R 0        | Q 0        | R 0        |            |            | B 0        |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| Q 0        |            |            | K 0        |            |            | Q 0        |            |
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
    white := 0
    black := 1

    for i := 0; i < t.N; i++ {
        b, err := newSimpleBoard(Point{8, 8}, 2)
        assert.Nil(t, err)

        setPiece(b, Point{0, 0}, Piece{black, QUEEN})
        setPiece(b, Point{3, 0}, Piece{black, KING_U})
        setPiece(b, Point{6, 0}, Piece{black, QUEEN})

        setPiece(b, Point{2, 1}, Piece{black, KNIGHT})
        setPiece(b, Point{3, 1}, Piece{black, QUEEN})
        setPiece(b, Point{4, 1}, Piece{black, KNIGHT})
        setPiece(b, Point{7, 1}, Piece{black, BISHOP})

        setPiece(b, Point{1, 2}, Piece{black, KNIGHT})
        setPiece(b, Point{3, 2}, Piece{white, PAWN_U})
        setPiece(b, Point{5, 2}, Piece{black, KNIGHT})
        setPiece(b, Point{7, 2}, Piece{black, BISHOP})

        setPiece(b, Point{2, 3}, Piece{black, ROOK})
        setPiece(b, Point{3, 3}, Piece{black, KNIGHT})
        setPiece(b, Point{4, 3}, Piece{black, ROOK})
        setPiece(b, Point{7, 3}, Piece{white, QUEEN})

        setPiece(b, Point{1, 4}, Piece{white, KNIGHT})
        setPiece(b, Point{2, 4}, Piece{white, QUEEN})
        setPiece(b, Point{4, 4}, Piece{white, QUEEN})
        setPiece(b, Point{5, 4}, Piece{white, KNIGHT})
        setPiece(b, Point{7, 4}, Piece{white, QUEEN})

        setPiece(b, Point{3, 5}, Piece{white, QUEEN})
        setPiece(b, Point{7, 5}, Piece{white, BISHOP})

        setPiece(b, Point{2, 6}, Piece{white, ROOK})
        setPiece(b, Point{3, 6}, Piece{white, QUEEN})
        setPiece(b, Point{4, 6}, Piece{white, ROOK})
        setPiece(b, Point{7, 6}, Piece{white, BISHOP})

        setPiece(b, Point{0, 7}, Piece{white, QUEEN})
        setPiece(b, Point{3, 7}, Piece{white, KING_U})
        setPiece(b, Point{6, 7}, Piece{white, QUEEN})

        // the time is around 1ms right now
        // when we skip the calculation it's around 40ns
        // when we skip the calculating checks, it's around 100us
        // b.test = true

        start := time.Now()
        b.CalculateMoves()
        end := time.Now()

        assert.Nil(t, err)
        fmt.Println("time", end.Sub(start))
    }
}

