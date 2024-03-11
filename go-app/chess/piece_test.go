package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPiece struct {
	mock.Mock
}

func (m *MockPiece) getColor() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockPiece) getValue() int {
    args := m.Called()
    return args.Int(0)
}

func (m *MockPiece) copy() (Piece, error) {
    args := m.Called()
    return args.Get(0).(Piece), args.Error(1)
}

func (m *MockPiece) getMoved() bool {
    args := m.Called()
    return args.Bool(0)
}

func (m *MockPiece) moves(board Board, location Point) []FastMove {
    args := m.Called(board, location)
	return args.Get(0).([]FastMove)
}

func (m *MockPiece) setDisabled(disabled bool) {
    m.Called(disabled)
}

func (m *MockPiece) getDisabled() bool {
    args := m.Called()
    return args.Bool(0)
}

func (m *MockPiece) print() string {
	args := m.Called()
	return args.String(0)
}

func Test_Pawn_Moves_Unmoved(t *testing.T) {
    white := 0
    pawn := pieceFactoryInstance.get(white, PAWN_D)

    b, err := newSimpleBoard(Point{7, 7}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{3, 3}, pawn)

	moves := pawn.moves(b, Point{3, 3})
    assert.Equal(t, 2, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{3, 4}])
    assert.True(t, moveMap[Point{3, 5}])
}

func Test_Pawn_Moves_Moved(t *testing.T) {
    white := 0
    pawn := pieceFactoryInstance.get(white, PAWN_D_M)

    b, err := newSimpleBoard(Point{7, 7}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{3, 3}, pawn)

	moves := pawn.moves(b, Point{3, 3})
    assert.Equal(t, 1, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{3, 4}])
}

func Test_Pawn_Moves_Capturing(t *testing.T) {
    white := 0
    black := 1
    pawn := pieceFactoryInstance.get(white, PAWN_D)
    blackPawn := pieceFactoryInstance.get(black, PAWN_D)

    b, err := newSimpleBoard(Point{7, 7}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{3, 3}, pawn)
    b.setPiece(Point{4, 4}, blackPawn)

	moves := pawn.moves(b, Point{3, 3})
    assert.Equal(t, 3, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{3, 4}])
    assert.True(t, moveMap[Point{3, 5}])
    assert.True(t, moveMap[Point{4, 4}])
}

func Test_Pawn_Moves_CapturingEnPassant(t *testing.T) {
    white := 0
    black := 1
    pawn := pieceFactoryInstance.get(white, PAWN_D)

    b, err := newSimpleBoard(Point{7, 7}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{3, 3}, pawn)
    b.setEnPassant(black, EnPassant{Point{4, 4}, Point{3, 4}})

	moves := pawn.moves(b, Point{3, 3})
    assert.Equal(t, 3, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{3, 4}])
    assert.True(t, moveMap[Point{3, 5}])
    assert.True(t, moveMap[Point{4, 4}])
}

func Test_Pawn_Moves_Promotion(t *testing.T) {
    white := 0
    pawn := pieceFactoryInstance.get(white, PAWN_D)

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{2, 2}, pawn)

	moves := pawn.moves(b, Point{2, 2})
    assert.Equal(t, 5, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
        moveMap[m.toLocation] = true
    }

    assert.True(t, moveMap[Point{2, 3}])
    assert.True(t, moveMap[Point{2, 4}])
}

func Test_Knight_Moves(t *testing.T) {
    white := 0
    knight := pieceFactoryInstance.get(white, KNIGHT)

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{2, 2}, knight)

	moves := knight.moves(b, Point{2, 2})
    assert.Equal(t, 8, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
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
    bishop := pieceFactoryInstance.get(white, BISHOP)

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{2, 2}, bishop)

	moves := bishop.moves(b, Point{2, 2})
    assert.Equal(t, 8, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
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
    rook := pieceFactoryInstance.get(white, ROOK)

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{2, 2}, rook)

	moves := rook.moves(b, Point{2, 2})
    assert.Equal(t, 8, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
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
    queen := pieceFactoryInstance.get(white, QUEEN)

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{2, 2}, queen)

	moves := queen.moves(b, Point{2, 2})
    assert.Equal(t, 16, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
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
    king := pieceFactoryInstance.get(white, KING_D)
    rook := pieceFactoryInstance.get(white, ROOK)

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{2, 2}, king)
    b.setPiece(Point{0, 2}, rook)
    b.setPiece(Point{4, 2}, rook)

	moves := king.moves(b, Point{2, 2})
    assert.Equal(t, 10, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
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
    king := pieceFactoryInstance.get(white, KING_D_M)
    rook := pieceFactoryInstance.get(white, ROOK)

    b, err := newSimpleBoard(Point{5, 5}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{2, 2}, king)
    b.setPiece(Point{0, 2}, rook)
    b.setPiece(Point{4, 2}, rook)

	moves := king.moves(b, Point{2, 2})
    assert.Equal(t, 8, len(moves))

    moveMap := map[Point]bool{}
    for _, m := range moves {
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

