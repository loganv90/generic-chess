package chess

import (
    "testing"
    "strings"
    "fmt"
    "time"

    "github.com/stretchr/testify/assert"
)

func Test_Minimax_depthOne(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(4, 4, 2)
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection(2)
    assert.Nil(t, err)

    i, err := invokerFactoryInstance.newSimpleInvoker()
    assert.Nil(t, err)

    b.setPiece(b.getIndex(2, 0), b.getAllPiece(black, KING_U))
    b.setPiece(b.getIndex(3, 0), b.getAllPiece(black, ROOK_M))
    b.setPiece(b.getIndex(2, 3), b.getAllPiece(white, KING_U))
    b.setPiece(b.getIndex(3, 3), b.getAllPiece(white, ROOK_M))
    b.CalculateMoves()

    game := &SimpleGame{
        b: b,
        p: p,
        i: i,
    }
    stop := make(chan bool)

    searcher := newSimpleSearcher(game, stop)
    moveKey, err := searcher.searchWithMinimax(1)
    assert.Nil(t, err)
    assert.Equal(t, 3, moveKey.XTo)
    assert.Equal(t, 0, moveKey.YTo)

    actualPrintedBoard := game.Print()
    expectedPrintedBoard := strings.Trim(`
+---------------------------------------------------+
|         0x |         1x |         2x |         3x |
|            |            | K 1        | R 1        |
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
|            |            | K 0        | R 0        |
|         3y |         3y |         3y |         3y |
+---------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

func Test_Minimax(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(4, 4, 2)
    assert.Nil(t, err)

    b.setPiece(b.getIndex(2, 3), b.getAllPiece(white, KING_U))
    b.setPiece(b.getIndex(3, 3), b.getAllPiece(white, ROOK))
    b.setPiece(b.getIndex(0, 0), b.getAllPiece(black, KING_D))
    b.disableLocation(b.getIndex(0, 3))
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
    stop := make(chan bool)

    searcher := newSimpleSearcher(game, stop)
    moveKey, err := searcher.searchWithMinimax(4)
    assert.Nil(t, err)
    assert.Equal(t, 1, moveKey.XTo)
    assert.Equal(t, 2, moveKey.YTo)

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

    b.setPiece(b.getIndex(5, 4), b.getAllPiece(white, QUEEN))
    b.setPiece(b.getIndex(5, 3), b.getAllPiece(white, QUEEN))
    b.setPiece(b.getIndex(3, 2), b.getAllPiece(black, PAWN_D_M))
    b.setPiece(b.getIndex(6, 0), nil)
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
    stop := make(chan bool)

    searcher := newSimpleSearcher(game, stop)
    moveKey, err := searcher.searchWithMinimax(3)
    assert.Nil(t, err)
    assert.Equal(t, 5, moveKey.XTo)
    assert.Equal(t, 2, moveKey.YTo)

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

// This is a problem with quiescence search and transposition maps
// We're evaluating positions where there's about to be a significant capture which makes the evaluation inaccurate
// Then we're storing that evaluation and retrieving it early in the search
// Evaluating the captures first seems to make this better but I don't think it solves it completely
func Test_Minimax_BotSacrifice(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(8, 8, 2)
    assert.Nil(t, err)

    b.setPiece(b.getIndex(1, 0), b.getAllPiece(black, ROOK))
    b.setPiece(b.getIndex(2, 0), b.getAllPiece(black, BISHOP))
    b.setPiece(b.getIndex(3, 0), b.getAllPiece(black, QUEEN))
    b.setPiece(b.getIndex(4, 0), b.getAllPiece(black, KING_D))
    b.setPiece(b.getIndex(7, 0), b.getAllPiece(black, ROOK))
    
    b.setPiece(b.getIndex(0, 1), b.getAllPiece(black, PAWN_D))
    b.setPiece(b.getIndex(2, 1), b.getAllPiece(black, PAWN_D))
    b.setPiece(b.getIndex(3, 1), b.getAllPiece(black, PAWN_D))
    b.setPiece(b.getIndex(4, 1), b.getAllPiece(black, KNIGHT))
    b.setPiece(b.getIndex(5, 1), b.getAllPiece(black, PAWN_D))
    b.setPiece(b.getIndex(6, 1), b.getAllPiece(black, PAWN_D))
    b.setPiece(b.getIndex(7, 1), b.getAllPiece(black, PAWN_D))

    b.setPiece(b.getIndex(4, 2), b.getAllPiece(black, PAWN_D_M))

    b.setPiece(b.getIndex(2, 3), b.getAllPiece(black, BISHOP))

    b.setPiece(b.getIndex(0, 4), b.getAllPiece(white, PAWN_U_M))
    b.setPiece(b.getIndex(4, 4), b.getAllPiece(white, PAWN_U_M))

    b.setPiece(b.getIndex(1, 5), b.getAllPiece(white, BISHOP))
    b.setPiece(b.getIndex(3, 5), b.getAllPiece(white, PAWN_U_M))
    b.setPiece(b.getIndex(5, 5), b.getAllPiece(white, KNIGHT))

    b.setPiece(b.getIndex(2, 6), b.getAllPiece(white, PAWN_U))
    b.setPiece(b.getIndex(5, 6), b.getAllPiece(white, PAWN_U))
    b.setPiece(b.getIndex(6, 6), b.getAllPiece(white, PAWN_U))
    b.setPiece(b.getIndex(7, 6), b.getAllPiece(white, PAWN_U))

    b.setPiece(b.getIndex(0, 7), b.getAllPiece(white, ROOK))
    b.setPiece(b.getIndex(2, 7), b.getAllPiece(white, BISHOP))
    b.setPiece(b.getIndex(3, 7), b.getAllPiece(white, QUEEN))
    b.setPiece(b.getIndex(4, 7), b.getAllPiece(white, KING_U))
    b.setPiece(b.getIndex(7, 7), b.getAllPiece(white, ROOK))

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
    stop := make(chan bool)

    searcher := newSimpleSearcher(game, stop)
    moveKey, err := searcher.searchWithMinimax(5)
    assert.Nil(t, err)
    assert.Equal(t, -1, moveKey.XTo)
    assert.Equal(t, -1, moveKey.YTo)

    actualPrintedBoard := game.Print()
    expectedPrintedBoard := strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            | R 1        | B 1        | Q 1        | K 1        |            |            | R 1        |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P 1        |            | P 1        | P 1        | N 1        | P 1        | P 1        | P 1        |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            | P 1        |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | B 1        |            |            |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P 0        |            |            |            | P 0        |            |            |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            | B 0        |            | P 0        |            | N 0        |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | P 0        |            |            | P 0        | P 0        | P 0        |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R 0        |            | B 0        | Q 0        | K 0        |            |            | R 0        |
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
        b, err := newSimpleBoard(8, 8, 2)
        assert.Nil(t, err)

        b.setPiece(b.getIndex(0, 0), b.getAllPiece(black, QUEEN))
        b.setPiece(b.getIndex(3, 0), b.getAllPiece(black, KING_U))
        b.setPiece(b.getIndex(6, 0), b.getAllPiece(black, QUEEN))

        b.setPiece(b.getIndex(2, 1), b.getAllPiece(black, KNIGHT))
        b.setPiece(b.getIndex(3, 1), b.getAllPiece(black, QUEEN))
        b.setPiece(b.getIndex(4, 1), b.getAllPiece(black, KNIGHT))
        b.setPiece(b.getIndex(7, 1), b.getAllPiece(black, BISHOP))

        b.setPiece(b.getIndex(1, 2), b.getAllPiece(black, KNIGHT))
        b.setPiece(b.getIndex(3, 2), b.getAllPiece(white, PAWN_U))
        b.setPiece(b.getIndex(5, 2), b.getAllPiece(black, KNIGHT))
        b.setPiece(b.getIndex(7, 2), b.getAllPiece(black, BISHOP))

        b.setPiece(b.getIndex(2, 3), b.getAllPiece(black, ROOK))
        b.setPiece(b.getIndex(3, 3), b.getAllPiece(black, KNIGHT))
        b.setPiece(b.getIndex(4, 3), b.getAllPiece(black, ROOK))
        b.setPiece(b.getIndex(7, 3), b.getAllPiece(white, QUEEN))

        b.setPiece(b.getIndex(1, 4), b.getAllPiece(white, KNIGHT))
        b.setPiece(b.getIndex(2, 4), b.getAllPiece(white, QUEEN))
        b.setPiece(b.getIndex(4, 4), b.getAllPiece(white, QUEEN))
        b.setPiece(b.getIndex(5, 4), b.getAllPiece(white, KNIGHT))
        b.setPiece(b.getIndex(7, 4), b.getAllPiece(white, QUEEN))

        b.setPiece(b.getIndex(3, 5), b.getAllPiece(white, QUEEN))
        b.setPiece(b.getIndex(7, 5), b.getAllPiece(white, BISHOP))

        b.setPiece(b.getIndex(2, 6), b.getAllPiece(white, ROOK))
        b.setPiece(b.getIndex(3, 6), b.getAllPiece(white, QUEEN))
        b.setPiece(b.getIndex(4, 6), b.getAllPiece(white, ROOK))
        b.setPiece(b.getIndex(7, 6), b.getAllPiece(white, BISHOP))

        b.setPiece(b.getIndex(0, 7), b.getAllPiece(white, QUEEN))
        b.setPiece(b.getIndex(3, 7), b.getAllPiece(white, KING_U))
        b.setPiece(b.getIndex(6, 7), b.getAllPiece(white, QUEEN))

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
        stop := make(chan bool)

        searcher := newSimpleSearcher(game, stop)
        _, err = searcher.searchWithMinimax(4)
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
    }
}

// without calculateMoves, the time falls to ~2s for depth 7
// without calculateToLocations, the time stays the same
// first we should try to calculate move for affected squared only then we can try more extreme solutions
func Benchmark_CalculateMoves(t *testing.B) {
    white := 0
    black := 1

    for i := 0; i < t.N; i++ {
        b, err := newSimpleBoard(8, 8, 2)
        assert.Nil(t, err)

        b.setPiece(b.getIndex(0, 0), b.getAllPiece(black, QUEEN))
        b.setPiece(b.getIndex(3, 0), b.getAllPiece(black, KING_U))
        b.setPiece(b.getIndex(6, 0), b.getAllPiece(black, QUEEN))

        b.setPiece(b.getIndex(2, 1), b.getAllPiece(black, KNIGHT))
        b.setPiece(b.getIndex(3, 1), b.getAllPiece(black, QUEEN))
        b.setPiece(b.getIndex(4, 1), b.getAllPiece(black, KNIGHT))
        b.setPiece(b.getIndex(7, 1), b.getAllPiece(black, BISHOP))

        b.setPiece(b.getIndex(1, 2), b.getAllPiece(black, KNIGHT))
        b.setPiece(b.getIndex(3, 2), b.getAllPiece(white, PAWN_U))
        b.setPiece(b.getIndex(5, 2), b.getAllPiece(black, KNIGHT))
        b.setPiece(b.getIndex(7, 2), b.getAllPiece(black, BISHOP))

        b.setPiece(b.getIndex(2, 3), b.getAllPiece(black, ROOK))
        b.setPiece(b.getIndex(3, 3), b.getAllPiece(black, KNIGHT))
        b.setPiece(b.getIndex(4, 3), b.getAllPiece(black, ROOK))
        b.setPiece(b.getIndex(7, 3), b.getAllPiece(white, QUEEN))

        b.setPiece(b.getIndex(1, 4), b.getAllPiece(white, KNIGHT))
        b.setPiece(b.getIndex(2, 4), b.getAllPiece(white, QUEEN))
        b.setPiece(b.getIndex(4, 4), b.getAllPiece(white, QUEEN))
        b.setPiece(b.getIndex(5, 4), b.getAllPiece(white, KNIGHT))
        b.setPiece(b.getIndex(7, 4), b.getAllPiece(white, QUEEN))

        b.setPiece(b.getIndex(3, 5), b.getAllPiece(white, QUEEN))
        b.setPiece(b.getIndex(7, 5), b.getAllPiece(white, BISHOP))

        b.setPiece(b.getIndex(2, 6), b.getAllPiece(white, ROOK))
        b.setPiece(b.getIndex(3, 6), b.getAllPiece(white, QUEEN))
        b.setPiece(b.getIndex(4, 6), b.getAllPiece(white, ROOK))
        b.setPiece(b.getIndex(7, 6), b.getAllPiece(white, BISHOP))

        b.setPiece(b.getIndex(0, 7), b.getAllPiece(white, QUEEN))
        b.setPiece(b.getIndex(3, 7), b.getAllPiece(white, KING_U))
        b.setPiece(b.getIndex(6, 7), b.getAllPiece(white, QUEEN))

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

