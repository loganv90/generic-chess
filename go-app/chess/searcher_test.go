package chess

import (
    "testing"
    "strings"

    "github.com/stretchr/testify/assert"
)

func Test_Minimax(t *testing.T) {
    b, err := newSimpleBoard(&Point{4, 4})
    assert.Nil(t, err)
    b.setPiece(&Point{2, 3}, newKing("white", false, 0, -1))
    b.setPiece(&Point{3, 3}, newRook("white", false))
    b.setPiece(&Point{0, 0}, newKing("black", false, 0, 1))
    b.disableLocation(&Point{0, 3})
    err = b.CalculateMoves()
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection([]*Player{{"white", true}, {"black", true}})
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

    score, moveKey, err := searcher.minimax(game, 3)
    assert.Nil(t, err)

    assert.Equal(t, 100000, score["white"])
    assert.Equal(t, -100000, score["black"])
    assert.Equal(t, 1, moveKey.XTo)
    assert.Equal(t, 2, moveKey.YTo)

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

