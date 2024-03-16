package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pawn_Moves_Unmoved(t *testing.T) {
    white := 0
    pawn := Piece{white, PAWN_D}

    b, err := newSimpleBoard(Point{7, 7}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{3, 3}, pawn)

    moves := Array100[FastMove]{}
	pawn.moves(b, &Point{3, 3}, &moves)
    assert.Equal(t, 2, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{3, 4}])
    assert.True(t, moveMap[Point{3, 5}])
}

func Test_Pawn_Moves_Moved(t *testing.T) {
    white := 0
    pawn := Piece{white, PAWN_D_M}

    b, err := newSimpleBoard(Point{7, 7}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{3, 3}, pawn)

    moves := Array100[FastMove]{}
	pawn.moves(b, &Point{3, 3}, &moves)
    assert.Equal(t, 1, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{3, 4}])
}

func Test_Pawn_Moves_Capturing(t *testing.T) {
    white := 0
    black := 1
    pawn := Piece{white, PAWN_D}
    blackPawn := Piece{black, PAWN_D}

    b, err := newSimpleBoard(Point{7, 7}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{3, 3}, pawn)
    setPiece(b, Point{4, 4}, blackPawn)

    moves := Array100[FastMove]{}
	pawn.moves(b, &Point{3, 3}, &moves)
    assert.Equal(t, 3, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{3, 4}])
    assert.True(t, moveMap[Point{3, 5}])
    assert.True(t, moveMap[Point{4, 4}])
}

func Test_Pawn_Moves_CapturingEnPassant(t *testing.T) {
    white := 0
    black := 1
    pawn := Piece{white, PAWN_D}

    b, err := newSimpleBoard(Point{7, 7}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{3, 3}, pawn)
    setEnPassant(b, black, EnPassant{Point{4, 4}, Point{3, 4}})

    moves := Array100[FastMove]{}
	pawn.moves(b, &Point{3, 3}, &moves)
    assert.Equal(t, 3, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{3, 4}])
    assert.True(t, moveMap[Point{3, 5}])
    assert.True(t, moveMap[Point{4, 4}])
}

func Test_Pawn_Moves_Promotion(t *testing.T) {
    white := 0
    pawn := Piece{white, PAWN_D}

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{2, 2}, pawn)

    moves := Array100[FastMove]{}
	pawn.moves(b, &Point{2, 2}, &moves)
    assert.Equal(t, 5, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{2, 3}])
    assert.True(t, moveMap[Point{2, 4}])
}

func Test_Knight_Moves(t *testing.T) {
    white := 0
    knight := Piece{white, KNIGHT}

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{2, 2}, knight)

    moves := Array100[FastMove]{}
	knight.moves(b, &Point{2, 2}, &moves)
    assert.Equal(t, 8, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{0, 1}])
    assert.True(t, moveMap[Point{1, 0}])
    assert.True(t, moveMap[Point{3, 0}])
    assert.True(t, moveMap[Point{4, 1}])
    assert.True(t, moveMap[Point{4, 3}])
    assert.True(t, moveMap[Point{3, 4}])
    assert.True(t, moveMap[Point{1, 4}])
    assert.True(t, moveMap[Point{0, 3}])
}

func Test_Bishop_Moves(t *testing.T) {
    white := 0
    bishop := Piece{white, BISHOP}

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{2, 2}, bishop)

    moves := Array100[FastMove]{}
	bishop.moves(b, &Point{2, 2}, &moves)
    assert.Equal(t, 8, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{0, 0}])
    assert.True(t, moveMap[Point{1, 1}])
    assert.True(t, moveMap[Point{3, 3}])
    assert.True(t, moveMap[Point{4, 4}])
    assert.True(t, moveMap[Point{4, 0}])
    assert.True(t, moveMap[Point{3, 1}])
    assert.True(t, moveMap[Point{1, 3}])
    assert.True(t, moveMap[Point{0, 4}])
}

func Test_Rook_Moves(t *testing.T) {
    white := 0
    rook := Piece{white, ROOK}

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{2, 2}, rook)

    moves := Array100[FastMove]{}
	rook.moves(b, &Point{2, 2}, &moves)
    assert.Equal(t, 8, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{2, 0}])
    assert.True(t, moveMap[Point{2, 1}])
    assert.True(t, moveMap[Point{2, 3}])
    assert.True(t, moveMap[Point{2, 4}])
    assert.True(t, moveMap[Point{0, 2}])
    assert.True(t, moveMap[Point{1, 2}])
    assert.True(t, moveMap[Point{3, 2}])
    assert.True(t, moveMap[Point{4, 2}])
}

func Test_Queen_Moves(t *testing.T) {
    white := 0
    queen := Piece{white, QUEEN}

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{2, 2}, queen)

    moves := Array100[FastMove]{}
	queen.moves(b, &Point{2, 2}, &moves)
    assert.Equal(t, 16, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{2, 0}])
    assert.True(t, moveMap[Point{2, 1}])
    assert.True(t, moveMap[Point{2, 3}])
    assert.True(t, moveMap[Point{2, 4}])
    assert.True(t, moveMap[Point{0, 2}])
    assert.True(t, moveMap[Point{1, 2}])
    assert.True(t, moveMap[Point{3, 2}])
    assert.True(t, moveMap[Point{4, 2}])
    assert.True(t, moveMap[Point{0, 0}])
    assert.True(t, moveMap[Point{1, 1}])
    assert.True(t, moveMap[Point{3, 3}])
    assert.True(t, moveMap[Point{4, 4}])
    assert.True(t, moveMap[Point{4, 0}])
    assert.True(t, moveMap[Point{3, 1}])
    assert.True(t, moveMap[Point{1, 3}])
    assert.True(t, moveMap[Point{0, 4}])
}

func Test_King_Moves_CanCastleAndUnmoved(t *testing.T) {
    white := 0
    king := Piece{white, KING_D}
    rook := Piece{white, ROOK}

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{2, 2}, king)
    setPiece(b, Point{0, 2}, rook)
    setPiece(b, Point{4, 2}, rook)

    moves := Array100[FastMove]{}
	king.moves(b, &Point{2, 2}, &moves)
    assert.Equal(t, 10, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{2, 1}])
    assert.True(t, moveMap[Point{2, 3}])
    assert.True(t, moveMap[Point{0, 2}])
    assert.True(t, moveMap[Point{1, 2}])
    assert.True(t, moveMap[Point{3, 2}])
    assert.True(t, moveMap[Point{4, 2}])
    assert.True(t, moveMap[Point{1, 1}])
    assert.True(t, moveMap[Point{3, 3}])
    assert.True(t, moveMap[Point{3, 1}])
    assert.True(t, moveMap[Point{1, 3}])
}

func Test_King_Moves_CanCastleAndMoved(t *testing.T) {
    white := 0
    king := Piece{white, KING_D_M}
    rook := Piece{white, ROOK}

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    setPiece(b, Point{2, 2}, king)
    setPiece(b, Point{0, 2}, rook)
    setPiece(b, Point{4, 2}, rook)

    moves := Array100[FastMove]{}
	king.moves(b, &Point{2, 2}, &moves)
    assert.Equal(t, 8, moves.count)

    moveMap := map[Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{2, 1}])
    assert.True(t, moveMap[Point{2, 3}])
    assert.True(t, moveMap[Point{1, 2}])
    assert.True(t, moveMap[Point{3, 2}])
    assert.True(t, moveMap[Point{1, 1}])
    assert.True(t, moveMap[Point{3, 3}])
    assert.True(t, moveMap[Point{3, 1}])
    assert.True(t, moveMap[Point{1, 3}])
}

