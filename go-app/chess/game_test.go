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
| R 1        | N 1        | B 1        | Q 1        | K 1        | B 1        | N 1        | R 1        |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P 1        | P 1        | P 1        | P 1        | P 1        | P 1        | P 1        | P 1        |
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
| P 0        | P 0        | P 0        | P 0        | P 0        | P 0        | P 0        | P 0        |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R 0        | N 0        | B 0        | Q 0        | K 0        | B 0        | N 0        | R 0        |
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
| R 1        | N 1        | B 1        | Q 1        | K 1        | B 1        | N 1        | R 1        |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P 1        | P 1        | P 1        | P 1        |            | P 1        | P 1        | P 1        |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            | P 1        |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            | P 0        |            |            |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P 0        | P 0        | P 0        | P 0        |            | P 0        | P 0        | P 0        |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R 0        | N 0        | B 0        | Q 0        | K 0        | B 0        | N 0        | R 0        |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

func Test_CastleWhenBlocked(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(Point{4, 4}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{1, 3}, pieceFactoryInstance.get(white, PAWN_U))
    b.setPiece(Point{2, 3}, pieceFactoryInstance.get(white, KING_U))
    b.setPiece(Point{3, 3}, pieceFactoryInstance.get(white, ROOK))
    b.setPiece(Point{0, 0}, pieceFactoryInstance.get(black, KING_D))
    err = b.CalculateMoves()
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection(2)
    assert.Nil(t, err)

    i, err := invokerFactoryInstance.newSimpleInvoker()
    assert.Nil(t, err)

    game := &SimpleGame{
        b: b,
        p: p,
        i: i,
    }

    err = game.Execute(2, 3, 3, 3, "")
    assert.NotNil(t, err)

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
|            | P 0        | K 0        | R 0        |
|         3y |         3y |         3y |         3y |
+---------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

func Test_CastleCanUndo(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(Point{4, 4}, 2)
    assert.Nil(t, err)
    b.setPiece(Point{2, 3}, pieceFactoryInstance.get(white, KING_U))
    b.setPiece(Point{3, 3}, pieceFactoryInstance.get(white, ROOK))
    b.setPiece(Point{0, 0}, pieceFactoryInstance.get(black, KING_D))
    err = b.CalculateMoves()
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection(2)
    assert.Nil(t, err)

    i, err := invokerFactoryInstance.newSimpleInvoker()
    assert.Nil(t, err)

    game := &SimpleGame{
        b: b,
        p: p,
        i: i,
    }

    err = game.Execute(2, 3, 3, 3, "")
    assert.Nil(t, err)

    err = game.Undo()
    assert.Nil(t, err)

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
|            |            | K 0        | R 0        |
|         3y |         3y |         3y |         3y |
+---------------------------------------------------+
	`, " \t\n") + "\n"
	assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

func Test_CapturePieceGivingCheck(t *testing.T) {
    white := 0

	game, err := NewSimpleFourPlayerGame()
	assert.Nil(t, err)

    err = game.Execute(7, 12, 7, 10, "") // white pawn advance
    assert.Nil(t, err)

    err = game.Execute(1, 7, 3, 7, "") // red pawn advance
    assert.Nil(t, err)

    err = game.Execute(7, 1, 7, 3, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Execute(12, 7, 10, 7, "") // blue pawn advance
    assert.Nil(t, err)

    err = game.Execute(6, 13, 13, 6, "") // white queen move
    assert.Nil(t, err)

    err = game.Execute(0, 6, 6, 12, "") // red queen move
    assert.Nil(t, err)

    err = game.Execute(6, 0, 12, 6, "") // black queen move
    assert.Nil(t, err)

    moves, err := game.Moves(white)
    assert.Nil(t, err)
    assert.Equal(t, 2, len(moves))
}

func Test_TwoPlayerCheckmate(t *testing.T) {
    white := 0

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
    assert.Equal(t, white, state.CurrentPlayer)
    assert.Equal(t, white, state.WinningPlayer)

    err = game.Undo()
    assert.Nil(t, err)

    state, err = game.State()
    assert.Nil(t, err)
    assert.Equal(t, white, state.CurrentPlayer)
    assert.Equal(t, -1, state.WinningPlayer)

    err = game.Redo()
    assert.Nil(t, err)

    moves, err := game.Moves(white)
    assert.Nil(t, err)
    err = game.Execute(moves[0].XFrom, moves[0].YFrom, moves[0].XTo, moves[0].YTo, moves[0].Promotion)
    assert.NotNil(t, err)
}

func Test_FourPlayerCheckmate(t *testing.T) {
    white := 0

    game, err := NewSimpleFourPlayerGame()
    assert.Nil(t, err)

    err = game.Execute(5, 12, 5, 11, "") // white pawn advance
    assert.Nil(t, err)
    err = game.Execute(1, 10, 2, 10, "") // red pawn advance
    assert.Nil(t, err)
    err = game.Execute(6, 1, 6, 2, "") // black pawn advance
    assert.Nil(t, err)
    err = game.Execute(12, 10, 11, 10, "") // blue pawn advance
    assert.Nil(t, err)

    err = game.Execute(7, 12, 7, 11, "") // white pawn advance
    assert.Nil(t, err)
    err = game.Execute(1, 3, 2, 3, "") // red pawn advance
    assert.Nil(t, err)
    err = game.Execute(5, 0, 7, 2, "") // black bishop advance
    assert.Nil(t, err)
    err = game.Execute(12, 3, 11, 3, "") // blue pawn advance
    assert.Nil(t, err)

    err = game.Execute(6, 13, 1, 8, "") // white queen checkmate
    assert.Nil(t, err)
    err = game.Execute(6, 2, 6, 3, "") // black pawn advance
    assert.Nil(t, err)
    err = game.Execute(12, 4, 11, 4, "") // blue pawn advance
    assert.Nil(t, err)

    err = game.Execute(8, 13, 2, 7, "") // white bishop advance
    assert.Nil(t, err)
    err = game.Execute(6, 0, 6, 2, "") // black queen advance
    assert.Nil(t, err)
    err = game.Execute(12, 5, 11, 5, "") // blue pawn advance
    assert.Nil(t, err)

    err = game.Execute(1, 8, 12, 8, "") // white queen checkmate
    assert.Nil(t, err)
    err = game.Execute(7, 2, 5, 0, "") // black bishop advance
    assert.Nil(t, err)

    err = game.Execute(12, 8, 8, 8, "") // white queen advance
    assert.Nil(t, err)
    err = game.Execute(6, 2, 6, 0, "") // black queen advance
    assert.Nil(t, err)

    err = game.Execute(2, 7, 6, 3, "") // white bishop advance
    assert.Nil(t, err)
    err = game.Execute(5, 0, 6, 1, "") // black bishop advance
    assert.Nil(t, err)

    err = game.Execute(6, 3, 8, 1, "") // white bishop checkmate
    assert.Nil(t, err)

    state, err := game.State()
    assert.Nil(t, err)
    assert.Equal(t, white, state.CurrentPlayer)
    assert.Equal(t, white, state.WinningPlayer)

    err = game.Undo()
    assert.Nil(t, err)

    state, err = game.State()
    assert.Nil(t, err)
    assert.Equal(t, white, state.CurrentPlayer)
    assert.Equal(t, -1, state.WinningPlayer)

    err = game.Redo()
    assert.Nil(t, err)

    moves, err := game.Moves(white)
    assert.Nil(t, err)
    err = game.Execute(moves[0].XFrom, moves[0].YFrom, moves[0].XTo, moves[0].YTo, moves[0].Promotion)
    assert.NotNil(t, err)
}

func Test_DisabledPieces(t *testing.T) {
    white := 0
    black := 1
    gray := 2

    b, err := newSimpleBoard(Point{8, 8}, 3)
    assert.Nil(t, err)
    b.setPiece(Point{0, 0}, pieceFactoryInstance.get(white, KING_D))
    b.setPiece(Point{7, 0}, pieceFactoryInstance.get(black, KING_D))
    b.setPiece(Point{4, 0}, pieceFactoryInstance.get(gray, KING_D))
    b.setPiece(Point{1, 7}, pieceFactoryInstance.get(black, QUEEN))
    b.setPiece(Point{6, 6}, pieceFactoryInstance.get(white, QUEEN))
    b.setPiece(Point{6, 7}, pieceFactoryInstance.get(white, QUEEN))
    err = b.CalculateMoves()
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection(3)
    assert.Nil(t, err)

    i, err := invokerFactoryInstance.newSimpleInvoker()
    assert.Nil(t, err)

    game := &SimpleGame{
        b: b,
        p: p,
        i: i,
    }

    err = game.Execute(6, 6, 6, 1, "") // white queen checkmate black
    assert.Nil(t, err)
    err = game.Execute(4, 0, 3, 0, "") // gray waiting move
    assert.Nil(t, err)
    err = game.Execute(0, 0, 1, 0, "") // white move into range of disabled black queen
    assert.Nil(t, err)

    state, err := game.State()
    assert.Nil(t, err)
    assert.Equal(t, -1, state.WinningPlayer)
    assert.Equal(t, false, state.GameOver)
}

func Test_Stalemate(t *testing.T) {
    white := 0
    black := 1
    gray := 2

    b, err := newSimpleBoard(Point{8, 8}, 3)
    assert.Nil(t, err)
    b.setPiece(Point{0, 0}, pieceFactoryInstance.get(white, KING_D))
    b.setPiece(Point{7, 0}, pieceFactoryInstance.get(black, KING_D))
    b.setPiece(Point{4, 0}, pieceFactoryInstance.get(gray, KING_D))
    b.setPiece(Point{6, 6}, pieceFactoryInstance.get(white, QUEEN))
    b.setPiece(Point{6, 7}, pieceFactoryInstance.get(white, QUEEN))
    err = b.CalculateMoves()
    assert.Nil(t, err)

    p, err := newSimplePlayerCollection(3)
    assert.Nil(t, err)

    i, err := invokerFactoryInstance.newSimpleInvoker()
    assert.Nil(t, err)

    game := &SimpleGame{
        b: b,
        p: p,
        i: i,
    }

    err = game.Execute(6, 6, 6, 2, "") // white queen stalemate black
    assert.Nil(t, err)

    state, err := game.State()
    assert.Nil(t, err)
    assert.Equal(t, -1, state.WinningPlayer)
    assert.Equal(t, true, state.GameOver)
}

func Test_NewSimpleGame(t *testing.T) {
	game, err := NewSimpleGame()
	assert.Nil(t, err)

	err = game.Execute(4, 6, 4, 4, "") // white pawn advance
	assert.Nil(t, err)

	err = game.Execute(1, 0, 2, 2, "") // black knight advance
	assert.Nil(t, err)

	err = game.Execute(4, 4, 4, 3, "") // white pawn advance
	assert.Nil(t, err)

	err = game.Execute(3, 1, 3, 3, "") // black pawn advance
	assert.Nil(t, err)

	err = game.Execute(4, 3, 3, 2, "") // white pawn capture en passant
	assert.Nil(t, err)

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

	err = game.Execute(4, 7, 7, 7, "") // white castle
	assert.Nil(t, err)

    err = game.Execute(0, 1, 0, 2, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Execute(3, 2, 4, 1, "") // white pawn capture
    assert.Nil(t, err)

    err = game.Execute(0, 2, 0, 3, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Execute(4, 1, 4, 0, "N") // white pawn capture
    assert.Nil(t, err)

    actualPrintedBoard := game.Print()
    expectedPrintedBoard := strings.Trim(`
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | K 1        | R 1        | N 0        | B 1        | N 1        | R 1        |
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            | P 1        | P 1        | Q 1        |            | P 1        | P 1        | P 1        |
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            | N 1        |            |            |            |            |            |
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P 1        | B 0        |            |            |            |            |            |            |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            |            | B 1        |            |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
|            |            |            |            |            | N 0        |            |            |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| P 0        | P 0        | P 0        | P 0        |            | P 0        | P 0        | P 0        |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |
| R 0        | N 0        | B 0        | Q 0        |            | R 0        | K 0        |            |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
    assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

func Test_NewSimpleFourPlayerGame(t *testing.T) {
	game, err := NewSimpleFourPlayerGame()
	assert.Nil(t, err)

    err = game.Execute(7, 12, 7, 10, "") // white pawn advance
    assert.Nil(t, err)

    err = game.Execute(1, 7, 3, 7, "") // red pawn advance
    assert.Nil(t, err)

    err = game.Execute(7, 1, 7, 3, "") // black pawn advance
    assert.Nil(t, err)

    err = game.Execute(12, 7, 10, 7, "") // blue pawn advance
    assert.Nil(t, err)

    err = game.Execute(8, 13, 6, 11, "") // white bishop move
    assert.Nil(t, err)

    err = game.Execute(0, 8, 2, 6, "") // red bishop move
    assert.Nil(t, err)

    err = game.Execute(8, 0, 6, 2, "") // black bishop move
    assert.Nil(t, err)

    err = game.Execute(13, 8, 11, 6, "") // blue bishop move
    assert.Nil(t, err)

    err = game.Execute(9, 13, 8, 11, "") // white knight move
    assert.Nil(t, err)

    err = game.Execute(0, 9, 2, 8, "") // red knight move
    assert.Nil(t, err)

    err = game.Execute(9, 0, 8, 2, "") // black knight move
    assert.Nil(t, err)

    err = game.Execute(13, 9, 11, 8, "") // blue knight move
    assert.Nil(t, err)

    err = game.Execute(7, 13, 10, 13, "") // white castle
    assert.Nil(t, err)

    err = game.Execute(0, 7, 0, 10, "") // red castle
    assert.Nil(t, err)

    err = game.Execute(7, 0, 10, 0, "") // black castle
    assert.Nil(t, err)

    err = game.Execute(13, 7, 13, 10, "") // blue castle
    assert.Nil(t, err)

    actualPrintedBoard := game.Print()
    expectedPrintedBoard := strings.Trim(`
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| R 2        | N 2        | B 2        | Q 2        |            | R 2        | K 2        |            |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |         0y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| P 2        | P 2        | P 2        | P 2        |            | P 2        | P 2        | P 2        |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |         1y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|            |            |            | B 2        |            | N 2        |            |            |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |         2y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| R 1        | P 1        |            |            |            |            |            | P 2        |            |            |            |            | P 3        | R 3        |
|         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |         3y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| N 1        | P 1        |            |            |            |            |            |            |            |            |            |            | P 3        | N 3        |
|         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |         4y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| B 1        | P 1        |            |            |            |            |            |            |            |            |            |            | P 3        | B 3        |
|         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |         5y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| Q 1        | P 1        | B 1        |            |            |            |            |            |            |            |            | B 3        | P 3        | Q 3        |
|         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |         6y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|            |            |            | P 1        |            |            |            |            |            |            | P 3        |            |            |            |
|         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |         7y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| R 1        | P 1        | N 1        |            |            |            |            |            |            |            |            | N 3        | P 3        | R 3        |
|         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |         8y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
| K 1        | P 1        |            |            |            |            |            |            |            |            |            |            | P 3        | K 3        |
|         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |         9y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|            | P 1        |            |            |            |            |            | P 0        |            |            |            |            | P 3        |            |
|        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |        10y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|            |            |            | B 0        |            | N 0        |            |            |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |        11y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| P 0        | P 0        | P 0        | P 0        |            | P 0        | P 0        | P 0        |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |        12y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|         0x |         1x |         2x |         3x |         4x |         5x |         6x |         7x |         8x |         9x |        10x |        11x |        12x |        13x |
|XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX| R 0        | N 0        | B 0        | Q 0        |            | R 0        | K 0        |            |XXXXXXXXXXXX|XXXXXXXXXXXX|XXXXXXXXXXXX|
|        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |        13y |
+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	`, " \t\n") + "\n"
    assert.Equal(t, expectedPrintedBoard, actualPrintedBoard)
}

