package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSimpleBoard_DefaultFen(t *testing.T) {
    b, err := createSimpleBoardWithDefaultPieceLocations()
	assert.Nil(t, err)

	for y := 2; y <= 5; y++ {
		for x := 0; x <= 7; x++ {
			piece := b.getPiece(b.getIndex(x, y))
			assert.True(t, piece == nil)
		}
	}

	for _, y := range []int{1, 6} {
		for x := 0; x <= 7; x++ {
            piece := b.getPiece(b.getIndex(x, y))
            assert.True(t, piece != nil)
            assert.True(t, piece.isPawn())
		}
	}

	for _, y := range []int{0, 7} {
        piece := b.getPiece(b.getIndex(0, y))
		assert.True(t, piece != nil)
        assert.True(t, piece.index == ROOK || piece.index == ROOK_M) // is rook

        piece = b.getPiece(b.getIndex(1, y))
		assert.True(t, piece != nil)
        assert.True(t, piece.index == KNIGHT) // is knight

        piece = b.getPiece(b.getIndex(2, y))
		assert.True(t, piece != nil)
        assert.True(t, piece.index == BISHOP) // is bishop

        piece = b.getPiece(b.getIndex(3, y))
		assert.True(t, piece != nil)
        assert.True(t, piece.index == QUEEN) // is queen

        piece = b.getPiece(b.getIndex(4, y))
		assert.True(t, piece != nil)
        assert.True(t, piece.isKing()) // is king

        piece = b.getPiece(b.getIndex(5, y))
		assert.True(t, piece != nil)
        assert.True(t, piece.index == BISHOP) // is bishop

        piece = b.getPiece(b.getIndex(6, y))
		assert.True(t, piece != nil)
        assert.True(t, piece.index == KNIGHT) // is knight

        piece = b.getPiece(b.getIndex(7, y))
		assert.True(t, piece != nil)
        assert.True(t, piece.index == ROOK || piece.index == ROOK_M) // is rook
	}
}

func Test_CalculateMoves_default(t *testing.T) {
    white := 0
    black := 1

    b, err := createSimpleBoardWithDefaultPieceLocations()
    assert.Nil(t, err)

    Assert_CountsAndMatest(t, b, white, 20, false, false, black, 20, false, false)
}

func Test_CalculateMoves_check(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(8, 8, 2)
    assert.Nil(t, err)

    b.setPiece(b.getIndex(0, 0), b.getAllPiece(white, KING_D))
    b.setPiece(b.getIndex(0, 1), b.getAllPiece(black, QUEEN))
    b.setPiece(b.getIndex(0, 7), b.getAllPiece(black, KING_U))

    Assert_CountsAndMatest(t, b, white, 1, false, false, black, 23, false, false)
}

func Test_CalculateMoves_checkmate(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(8, 8, 2)
    assert.Nil(t, err)

    b.setPiece(b.getIndex(0, 0), b.getAllPiece(white, KING_D))
    b.setPiece(b.getIndex(0, 1), b.getAllPiece(black, QUEEN))
    b.setPiece(b.getIndex(0, 2), b.getAllPiece(black, KING_U))

    Assert_CountsAndMatest(t, b, white, 0, true, false, black, 18, false, false)
}

func Test_CalculateMoves_stalemate(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(8, 8, 2)
    assert.Nil(t, err)

    b.setPiece(b.getIndex(0, 0), b.getAllPiece(white, KING_D))
    b.setPiece(b.getIndex(1, 2), b.getAllPiece(black, QUEEN))
    b.setPiece(b.getIndex(0, 7), b.getAllPiece(black, KING_U))

    Assert_CountsAndMatest(t, b, white, 0, false, true, black, 26, false, false)
}

func Test_CalculateMoves_noCastleThroughCheck(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(8, 8, 2)
    assert.Nil(t, err)

    b.setPiece(b.getIndex(4, 0), b.getAllPiece(white, KING_D))
    b.setPiece(b.getIndex(0, 0), b.getAllPiece(white, ROOK_M))
    b.setPiece(b.getIndex(4, 7), b.getAllPiece(black, KING_U))
    b.setPiece(b.getIndex(3, 7), b.getAllPiece(black, ROOK_M))

    Assert_CountsAndMatest(t, b, white, 13, false, false, black, 14, false, false)
}

func Test_CalculateMoves_castle(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(8, 8, 2)
    assert.Nil(t, err)

    b.setPiece(b.getIndex(4, 0), b.getAllPiece(white, KING_D))
    b.setPiece(b.getIndex(0, 0), b.getAllPiece(white, ROOK))
    b.setPiece(b.getIndex(4, 7), b.getAllPiece(black, KING_U))

    Assert_CountsAndMatest(t, b, white, 16, false, false, black, 5, false, false)
}

func Test_CalculateMoves_promotion(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(8, 8, 2)
    assert.Nil(t, err)

    b.setPiece(b.getIndex(7, 6), b.getAllPiece(white, PAWN_D_M))
    b.setPiece(b.getIndex(0, 0), b.getAllPiece(white, KING_D))
    b.setPiece(b.getIndex(0, 7), b.getAllPiece(black, KING_U))

    Assert_CountsAndMatest(t, b, white, 7, false, false, black, 3, false, false)
}

func Assert_CountsAndMatest(
    t *testing.T,
    b Board,
    white int,
    whiteCount int,
    whiteCheckmate bool,
    whiteStalemate bool,
    black int,
    blackCount int,
    blackCheckmate bool,
    blackStalemate bool,
) {
    b.CalculateMoves()

    whiteMoveKeys, err := b.LegalMovesOfColor(white)
    assert.Nil(t, err)
    assert.Equal(t, whiteCount, len(whiteMoveKeys))

    blackMoveKeys, err := b.LegalMovesOfColor(black)
    assert.Nil(t, err)
    assert.Equal(t, blackCount, len(blackMoveKeys))

    checkmate, stalemate, err := b.CheckmateAndStalemate(white)
    assert.Nil(t, err)
    assert.Equal(t, whiteCheckmate, checkmate)
    assert.Equal(t, whiteStalemate, stalemate)

    checkmate, stalemate, err = b.CheckmateAndStalemate(black)
    assert.Nil(t, err)
    assert.Equal(t, blackCheckmate, checkmate)
    assert.Equal(t, blackStalemate, stalemate)
}

// TODO add en passant stuff to the unique string
func Test_MinimumString_Default(t *testing.T) {
    b, err := createSimpleBoardWithDefaultPieceLocations()
    assert.Nil(t, err)

    expected := "R1N1B1Q1K1B1N1R1P1P1P1P1P1P1P1P132P0P0P0P0P0P0P0P0R0N0B0Q0K0B0N0R0"
    actual := b.UniqueString()
    assert.Equal(t, expected, actual)
}

