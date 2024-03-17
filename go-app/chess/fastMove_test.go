package chess

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_MoveSimple(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    b.setEnPassant(white, b.getIndex(2, 2), b.getIndex(3, 3))
    b.setVulnerable(white, b.getIndex(4, 4), b.getIndex(5, 5))

    moves := Array1000[FastMove]{}
    move := moves.get()
    addMoveSimple(
        b,
        b.getAllPiece(white, PAWN_D),
        b.getIndex(0, 0),
        b.getAllPiece(black, PAWN_D),
        b.getIndex(1, 1),
        nil,
        &moves,
    )
    assert.Equal(t, move.allyDefense, false)

    move.execute()
    assert.True(t, b.getPiece(b.getIndex(0, 0)) == nil)
    assert.Equal(t, b.getPiece(b.getIndex(1, 1)), b.getAllPiece(white, PAWN_D_M))
    target, risk := b.getEnPassant(white)
    start, end := b.getVulnerable(white)
    assert.True(t, target == nil)
    assert.True(t, risk == nil)
    assert.True(t, start == nil)
    assert.True(t, end == nil)

    move.undo()
    assert.Equal(t, b.getPiece(b.getIndex(0, 0)), b.getAllPiece(white, PAWN_D))
    assert.Equal(t, b.getPiece(b.getIndex(1, 1)), b.getAllPiece(black, PAWN_D))
    target, risk = b.getEnPassant(white)
    start, end = b.getVulnerable(white)
    assert.Equal(t, target, b.getIndex(2, 2))
    assert.Equal(t, risk, b.getIndex(3, 3))
    assert.Equal(t, start, b.getIndex(4, 4))
    assert.Equal(t, end, b.getIndex(5, 5))
}

func Test_MovePromotion(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    b.setEnPassant(white, b.getIndex(2, 2), b.getIndex(3, 3))
    b.setVulnerable(white, b.getIndex(4, 4), b.getIndex(5, 5))

    moves := Array1000[FastMove]{}
    move := moves.get()
    addMoveSimple(
        b,
        b.getAllPiece(white, PAWN_D),
        b.getIndex(0, 0),
        b.getAllPiece(black, PAWN_D),
        b.getIndex(1, 1),
        b.getAllPiece(white, QUEEN),
        &moves,
    )
    assert.Equal(t, move.allyDefense, false)

    move.execute()
    assert.True(t, b.getPiece(b.getIndex(0, 0)) == nil)
    assert.Equal(t, b.getPiece(b.getIndex(1, 1)), b.getAllPiece(white, QUEEN))
    target, risk := b.getEnPassant(white)
    start, end := b.getVulnerable(white)
    assert.True(t, target == nil)
    assert.True(t, risk == nil)
    assert.True(t, start == nil)
    assert.True(t, end == nil)

    move.undo()
    assert.Equal(t, b.getPiece(b.getIndex(0, 0)), b.getAllPiece(white, PAWN_D))
    assert.Equal(t, b.getPiece(b.getIndex(1, 1)), b.getAllPiece(black, PAWN_D))
    target, risk = b.getEnPassant(white)
    start, end = b.getVulnerable(white)
    assert.Equal(t, target, b.getIndex(2, 2))
    assert.Equal(t, risk, b.getIndex(3, 3))
    assert.Equal(t, start, b.getIndex(4, 4))
    assert.Equal(t, end, b.getIndex(5, 5))
}

func Test_MoveRevealEnPassant(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    b.setEnPassant(white, b.getIndex(2, 2), b.getIndex(3, 3))
    b.setVulnerable(white, b.getIndex(4, 4), b.getIndex(5, 5))

    moves := Array1000[FastMove]{}
    move := moves.get()
    addMoveRevealEnPassant(
        b,
        b.getAllPiece(white, PAWN_D),
        b.getIndex(0, 0),
        b.getAllPiece(black, PAWN_D),
        b.getIndex(1, 1),
        nil,
        b.getIndex(6, 6),
        b.getIndex(7, 7),
        &moves,
    )
    assert.Equal(t, move.allyDefense, false)

    move.execute()
    assert.True(t, b.getPiece(b.getIndex(0, 0)) == nil)
    assert.Equal(t, b.getPiece(b.getIndex(1, 1)), b.getAllPiece(white, PAWN_D_M))
    target, risk := b.getEnPassant(white)
    start, end := b.getVulnerable(white)
    assert.Equal(t, target, b.getIndex(6, 6))
    assert.Equal(t, risk, b.getIndex(7, 7))
    assert.True(t, start == nil)
    assert.True(t, end == nil)

    move.undo()
    assert.Equal(t, b.getPiece(b.getIndex(0, 0)), b.getAllPiece(white, PAWN_D))
    assert.Equal(t, b.getPiece(b.getIndex(1, 1)), b.getAllPiece(black, PAWN_D))
    target, risk = b.getEnPassant(white)
    start, end = b.getVulnerable(white)
    assert.Equal(t, target, b.getIndex(2, 2))
    assert.Equal(t, risk, b.getIndex(3, 3))
    assert.Equal(t, start, b.getIndex(4, 4))
    assert.Equal(t, end, b.getIndex(5, 5))
}

func Test_MoveCaptureEnPassant(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    b.setPiece(b.getIndex(8, 8), b.getAllPiece(black, PAWN_D))
    b.setPiece(b.getIndex(9, 9), b.getAllPiece(black, PAWN_D))
    b.setEnPassant(white, b.getIndex(2, 2), b.getIndex(3, 3))
    b.setVulnerable(white, b.getIndex(4, 4), b.getIndex(5, 5))

    moves := Array1000[FastMove]{}
    move := moves.get()
    addMoveCaptureEnPassant(
        b,
        b.getAllPiece(white, PAWN_D),
        b.getIndex(0, 0),
        b.getAllPiece(black, PAWN_D),
        b.getIndex(1, 1),
        nil,
        b.getIndex(6, 6),
        b.getIndex(7, 7),
        b.getIndex(8, 8),
        b.getIndex(9, 9),
        &moves,
    )
    assert.Equal(t, move.allyDefense, false)

    move.execute()
    assert.True(t, b.getPiece(b.getIndex(0, 0)) == nil)
    assert.Equal(t, b.getPiece(b.getIndex(1, 1)), b.getAllPiece(white, PAWN_D_M))
    assert.True(t, b.getPiece(b.getIndex(8, 8)) == nil)
    assert.True(t, b.getPiece(b.getIndex(9, 9)) == nil)
    target, risk := b.getEnPassant(white)
    start, end := b.getVulnerable(white)
    assert.True(t, target == nil)
    assert.True(t, risk == nil)
    assert.True(t, start == nil)
    assert.True(t, end == nil)

    move.undo()
    assert.Equal(t, b.getPiece(b.getIndex(0, 0)), b.getAllPiece(white, PAWN_D))
    assert.Equal(t, b.getPiece(b.getIndex(1, 1)), b.getAllPiece(black, PAWN_D))
    assert.Equal(t, b.getPiece(b.getIndex(8, 8)), b.getAllPiece(black, PAWN_D))
    assert.Equal(t, b.getPiece(b.getIndex(9, 9)), b.getAllPiece(black, PAWN_D))
    target, risk = b.getEnPassant(white)
    start, end = b.getVulnerable(white)
    assert.Equal(t, target, b.getIndex(2, 2))
    assert.Equal(t, risk, b.getIndex(3, 3))
    assert.Equal(t, start, b.getIndex(4, 4))
    assert.Equal(t, end, b.getIndex(5, 5))
}

func Test_MoveAllyDefense(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    b.setEnPassant(white, b.getIndex(2, 2), b.getIndex(3, 3))
    b.setVulnerable(white, b.getIndex(4, 4), b.getIndex(5, 5))

    moves := Array1000[FastMove]{}
    move := moves.get()
    addMoveAllyDefense(
        b,
        b.getAllPiece(white, PAWN_D),
        b.getIndex(0, 0),
        b.getIndex(1, 1),
        &moves,
    )
    assert.Equal(t, move.allyDefense, true)

    move.execute()
    assert.True(t, b.getPiece(b.getIndex(0, 0)) == nil)
    assert.True(t, b.getPiece(b.getIndex(1, 1)) == nil)
    target, risk := b.getEnPassant(white)
    start, end := b.getVulnerable(white)
    assert.True(t, target == nil)
    assert.True(t, risk == nil)
    assert.True(t, start == nil)
    assert.True(t, end == nil)

    move.undo()
    assert.True(t, b.getPiece(b.getIndex(0, 0)) == nil)
    assert.True(t, b.getPiece(b.getIndex(1, 1)) == nil)
    target, risk = b.getEnPassant(white)
    start, end = b.getVulnerable(white)
    assert.True(t, target == nil)
    assert.True(t, risk == nil)
    assert.True(t, start == nil)
    assert.True(t, end == nil)
}

func Test_MoveCastle(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(10, 10, 2)
    assert.Nil(t, err)

    b.setEnPassant(white, b.getIndex(2, 2), b.getIndex(3, 3))
    b.setVulnerable(white, b.getIndex(4, 4), b.getIndex(5, 5))

    moves := Array1000[FastMove]{}
    move := moves.get()
    addMoveCastle(
        b,
        b.getAllPiece(white, KING_D),
        b.getIndex(0, 0),
        b.getIndex(6, 6),
        b.getAllPiece(white, ROOK),
        b.getIndex(1, 1),
        b.getIndex(7, 7),
        b.getIndex(8, 8),
        b.getIndex(9, 9),
        &moves,
    )
    assert.Equal(t, move.allyDefense, false)

    move.execute()
    assert.True(t, b.getPiece(b.getIndex(0, 0)) == nil)
    assert.True(t, b.getPiece(b.getIndex(1, 1)) == nil)
    assert.Equal(t, b.getPiece(b.getIndex(6, 6)), b.getAllPiece(white, KING_D_M))
    assert.Equal(t, b.getPiece(b.getIndex(7, 7)), b.getAllPiece(white, ROOK_M))
    target, risk := b.getEnPassant(white)
    start, end := b.getVulnerable(white)
    assert.True(t, target == nil)
    assert.True(t, risk == nil)
    assert.Equal(t, start, b.getIndex(8, 8))
    assert.Equal(t, end, b.getIndex(9, 9))

    move.undo()
    assert.Equal(t, b.getPiece(b.getIndex(0, 0)), b.getAllPiece(white, KING_D))
    assert.Equal(t, b.getPiece(b.getIndex(1, 1)), b.getAllPiece(white, ROOK))
    assert.True(t, b.getPiece(b.getIndex(6, 6)) == nil)
    assert.True(t, b.getPiece(b.getIndex(7, 7)) == nil)
    target, risk = b.getEnPassant(white)
    start, end = b.getVulnerable(white)
    assert.Equal(t, target, b.getIndex(2, 2))
    assert.Equal(t, risk, b.getIndex(3, 3))
    assert.Equal(t, start, b.getIndex(4, 4))
    assert.Equal(t, end, b.getIndex(5, 5))
}

