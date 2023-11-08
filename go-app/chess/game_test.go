package chess

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSimpleGame_EnPassant(t *testing.T) {
	game, err := NewSimpleGame()
	assert.Nil(t, err)

	actualPrintedBoard := game.print()
	expectedPrintedBoard := strings.Trim(`
Player: white
Check:  false
Mate:   false
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
	assert.Equal(t, actualPrintedBoard, expectedPrintedBoard)

	_, err = game.Execute(4, 6, 4, 4) // white pawn advance
	assert.Nil(t, err)

	actualPrintedBoard = game.print()
	expectedPrintedBoard = strings.Trim(`
Player: black
Check:  false
Mate:   false
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
	assert.Equal(t, actualPrintedBoard, expectedPrintedBoard)

	_, err = game.Execute(1, 0, 2, 2) // black knight advance
	assert.Nil(t, err)

	actualPrintedBoard = game.print()
	expectedPrintedBoard = strings.Trim(`
Player: white
Check:  false
Mate:   false
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
	assert.Equal(t, actualPrintedBoard, expectedPrintedBoard)

	_, err = game.Execute(4, 4, 4, 3) // white pawn advance
	assert.Nil(t, err)

	_, err = game.Execute(3, 1, 3, 3) // black pawn advance
	assert.Nil(t, err)

	_, err = game.Execute(4, 3, 3, 2) // white pawn capture en passant
	assert.Nil(t, err)

	actualPrintedBoard = game.print()
	expectedPrintedBoard = strings.Trim(`
Player: black
Check:  false
Mate:   false
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
	assert.Equal(t, actualPrintedBoard, expectedPrintedBoard)

	_, err = game.Execute(2, 0, 6, 4) // black bishop advance
	assert.Nil(t, err)

	_, err = game.Execute(5, 7, 1, 3) // white bishop advance
	assert.Nil(t, err)

	_, err = game.Execute(3, 0, 3, 1) // white queen advance
	assert.Nil(t, err)

	_, err = game.Execute(6, 7, 5, 5) // white knight advance
	assert.Nil(t, err)

	_, err = game.Execute(4, 0, 0, 0) // black castle
	assert.Nil(t, err)

	actualPrintedBoard = game.print()
	expectedPrintedBoard = strings.Trim(`
Player: white
Check:  false
Mate:   false
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
	assert.Equal(t, actualPrintedBoard, expectedPrintedBoard)

	_, err = game.Execute(4, 7, 7, 7) // white castle
	assert.Nil(t, err)

	actualPrintedBoard = game.print()
	expectedPrintedBoard = strings.Trim(`
Player: black
Check:  false
Mate:   false
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
	assert.Equal(t, actualPrintedBoard, expectedPrintedBoard)
}
