package chess

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UndoAndRedo(t *testing.T) {
    game, err := NewSimpleGame()
    assert.Nil(t, err)

    err = game.Execute(4, 6, 4, 4, "") // white pawn advance
    assert.Nil(t, err)

    err = game.Execute(4, 1, 4, 3, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Undo()
    assert.Nil(t, err)

    err = game.Undo()
    assert.Nil(t, err)

    err = game.Undo()
    assert.NotNil(t, err)

	actualPrintedBoard := game.Print()
	expectedPrintedBoard := strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R black    | N black    | B black    | Q black    | K black    | B black    | N black    | R black    |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P black    | P black    | P black    | P black    | P black    | P black    | P black    | P black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
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

    err = game.Redo()
    assert.Nil(t, err)

    err = game.Redo()
    assert.Nil(t, err)

    err = game.Redo()
    assert.NotNil(t, err)

	actualPrintedBoard = game.Print()
	expectedPrintedBoard = strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R black    | N black    | B black    | Q black    | K black    | B black    | N black    | R black    |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P black    | P black    | P black    | P black    |            | P black    | P black    | P black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            | P black    |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            | P white    |            |            |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P white    | P white    | P white    | P white    |            | P white    | P white    | P white    |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R white    | N white    | B white    | Q white    | K white    | B white    | N white    | R white    |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

// TODO add four player checkmate test
// TODO add some kind of event system to the state to notify checkmates and stalemates
func Test_TwoPlayerCheckmate(t *testing.T) {
    game, err := NewSimpleGame()
    assert.Nil(t, err)

    err = game.Execute(4, 6, 4, 4, "") // white pawn advance
    assert.Nil(t, err)

    err = game.Execute(0, 1, 0, 2, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Execute(5, 7, 2, 4, "") // white bishop advance
    assert.Nil(t, err)

    err = game.Execute(0, 2, 0, 3, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Execute(3, 7, 5, 5, "") // white queen advance
    assert.Nil(t, err)

    err = game.Execute(0, 3, 0, 4, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Execute(5, 5, 5, 1, "") // white queen checkmate
    assert.Nil(t, err)

    state, err := game.State()
    assert.Nil(t, err)
    assert.Equal(t, "white", state.CurrentPlayer)
    assert.Equal(t, "white", state.WinningPlayer)

    err = game.Undo()
    assert.Nil(t, err)

    state, err = game.State()
    assert.Nil(t, err)
    assert.Equal(t, "white", state.CurrentPlayer)
    assert.Equal(t, "", state.WinningPlayer)
}

func Test_NewSimpleGame(t *testing.T) {
	game, err := NewSimpleGame()
	assert.Nil(t, err)

	actualPrintedBoard := game.Print()
	expectedPrintedBoard := strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R black    | N black    | B black    | Q black    | K black    | B black    | N black    | R black    |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P black    | P black    | P black    | P black    | P black    | P black    | P black    | P black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
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

	err = game.Execute(4, 6, 4, 4, "") // white pawn advance
	assert.Nil(t, err)

	actualPrintedBoard = game.Print()
	expectedPrintedBoard = strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R black    | N black    | B black    | Q black    | K black    | B black    | N black    | R black    |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P black    | P black    | P black    | P black    | P black    | P black    | P black    | P black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            | P white    |            |            |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P white    | P white    | P white    | P white    |            | P white    | P white    | P white    |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R white    | N white    | B white    | Q white    | K white    | B white    | N white    | R white    |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)

	err = game.Execute(1, 0, 2, 2, "") // black knight advance
	assert.Nil(t, err)

	actualPrintedBoard = game.Print()
	expectedPrintedBoard = strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R black    |            | B black    | Q black    | K black    | B black    | N black    | R black    |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P black    | P black    | P black    | P black    | P black    | P black    | P black    | P black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | N black    |            |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            | P white    |            |            |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P white    | P white    | P white    | P white    |            | P white    | P white    | P white    |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R white    | N white    | B white    | Q white    | K white    | B white    | N white    | R white    |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)

	err = game.Execute(4, 4, 4, 3, "") // white pawn advance
	assert.Nil(t, err)

	err = game.Execute(3, 1, 3, 3, "") // black pawn advance
	assert.Nil(t, err)

	err = game.Execute(4, 3, 3, 2, "") // white pawn capture en passant
	assert.Nil(t, err)

	actualPrintedBoard = game.Print()
	expectedPrintedBoard = strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R black    |            | B black    | Q black    | K black    | B black    | N black    | R black    |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P black    | P black    | P black    |            | P black    | P black    | P black    | P black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | N black    | P white    |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P white    | P white    | P white    | P white    |            | P white    | P white    | P white    |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R white    | N white    | B white    | Q white    | K white    | B white    | N white    | R white    |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)

	err = game.Execute(2, 0, 6, 4, "") // black bishop advance
	assert.Nil(t, err)

	err = game.Execute(5, 7, 1, 3, "") // white bishop advance
	assert.Nil(t, err)

	err = game.Execute(3, 0, 3, 1, "") // white queen advance
	assert.Nil(t, err)

	err = game.Execute(6, 7, 5, 5, "") // white knight advance
	assert.Nil(t, err)

	err = game.Execute(4, 0, 0, 0, "") // black castle
	assert.Nil(t, err)

	actualPrintedBoard = game.Print()
	expectedPrintedBoard = strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | K black    | R black    |            | B black    | N black    | R black    |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P black    | P black    | P black    | Q black    | P black    | P black    | P black    | P black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | N black    | P white    |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            | B white    |            |            |            |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            | B black    |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            | N white    |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P white    | P white    | P white    | P white    |            | P white    | P white    | P white    |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R white    | N white    | B white    | Q white    | K white    |            |            | R white    |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)

	err = game.Execute(4, 7, 7, 7, "") // white castle
	assert.Nil(t, err)

	actualPrintedBoard = game.Print()
	expectedPrintedBoard = strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | K black    | R black    |            | B black    | N black    | R black    |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P black    | P black    | P black    | Q black    | P black    | P black    | P black    | P black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | N black    | P white    |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            | B white    |            |            |            |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            | B black    |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            | N white    |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P white    | P white    | P white    | P white    |            | P white    | P white    | P white    |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R white    | N white    | B white    | Q white    |            | R white    | K white    |            |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)

    err = game.Execute(0, 1, 0, 2, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Execute(3, 2, 4, 1, "") // white pawn capture
    assert.Nil(t, err)

    err = game.Execute(0, 2, 0, 3, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Execute(4, 1, 4, 0, "N") // white pawn capture
    assert.Nil(t, err)

    actualPrintedBoard = game.Print()
    expectedPrintedBoard = strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | K black    | R black    | N white    | B black    | N black    | R black    |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            | P black    | P black    | Q black    |            | P black    | P black    | P black    |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | N black    |            |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P black    | B white    |            |            |            |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            | B black    |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            | N white    |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P white    | P white    | P white    | P white    |            | P white    | P white    | P white    |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R white    | N white    | B white    | Q white    |            | R white    | K white    |            |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
    assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

func Test_NewSimpleFourPlayerGame(t *testing.T) {
	game, err := NewSimpleFourPlayerGame()
	assert.Nil(t, err)

	actualPrintedBoard := game.Print()
	expectedPrintedBoard := strings.Trim(`
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| R black    | N black    | B black    | Q black    | K black    | B black    | N black    | R black    |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| P black    | P black    | P black    | P black    | P black    | P black    | P black    | P black    |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|            |            |            |            |            |            |            |            |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| R red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | R blue     |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| N red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | N blue     |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| B red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | B blue     |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| Q red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | Q blue     |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| K red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | K blue     |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| B red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | B blue     |
|         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| N red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | N blue     |
|         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| R red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | R blue     |
|        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|            |            |            |            |            |            |            |            |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| P white    | P white    | P white    | P white    | P white    | P white    | P white    | P white    |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| R white    | N white    | B white    | Q white    | K white    | B white    | N white    | R white    |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)

    err = game.Execute(7, 12, 7, 10, "") // white pawn advance
    assert.Nil(t, err)

    err = game.Execute(1, 7, 3, 7, "") // red pawn advance
    assert.Nil(t, err)

    err = game.Execute(7, 1, 7, 3, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Execute(12, 7, 10, 7, "") // blue pawn advance
    assert.Nil(t, err)

    actualPrintedBoard = game.Print()
    expectedPrintedBoard = strings.Trim(`
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| R black    | N black    | B black    | Q black    | K black    | B black    | N black    | R black    |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| P black    | P black    | P black    | P black    |            | P black    | P black    | P black    |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|            |            |            |            |            |            |            |            |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| R red      | P red      |            |            |            |            |            | P black    |            |            |            |            | P blue     | R blue     |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| N red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | N blue     |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| B red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | B blue     |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| Q red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | Q blue     |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| K red      |            |            | P red      |            |            |            |            |            |            | P blue     |            |            | K blue     |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| B red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | B blue     |
|         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| N red      | P red      |            |            |            |            |            |            |            |            |            |            | P blue     | N blue     |
|         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| R red      | P red      |            |            |            |            |            | P white    |            |            |            |            | P blue     | R blue     |
|        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|            |            |            |            |            |            |            |            |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| P white    | P white    | P white    | P white    |            | P white    | P white    | P white    |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| R white    | N white    | B white    | Q white    | K white    | B white    | N white    | R white    |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
    assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

