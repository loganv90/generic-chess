package chess


import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBoard struct {
	mock.Mock
}

func (m *MockBoard) getPiece(location *Point) (Piece, error) {
	args := m.Called(location)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(Piece), args.Error(1)
	}
}

func (m *MockBoard) setPiece(location *Point, piece Piece) error {
	args := m.Called(location, piece)
	return args.Error(0)
}

func (m *MockBoard) getVulnerables(color string) ([]*Point, error) {
    args := m.Called(color)
    return args.Get(0).([]*Point), args.Error(1)
}

func (m *MockBoard) setVulnerables(color string, vulnerables []*Point) error {
    args := m.Called(color, vulnerables)
    return args.Error(0)
}

func (m *MockBoard) getEnPassant(color string) (*EnPassant, error) {
	args := m.Called(color)
	return args.Get(0).(*EnPassant), args.Error(1)
}

func (m *MockBoard) setEnPassant(color string, enPassant *EnPassant) error {
    args := m.Called(color, enPassant)
    return args.Error(0)
}

func (m *MockBoard) possibleEnPassant(color string, location *Point) ([]*EnPassant, error) {
    args := m.Called(color, location)
    return args.Get(0).([]*EnPassant), args.Error(1)
}

func (m *MockBoard) clearEnPassant(color string) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockBoard) PotentialMoves(fromLocation *Point) ([]Move, error) {
    args := m.Called(fromLocation)
    return args.Get(0).([]Move), args.Error(1)
}

func (m *MockBoard) ValidMoves(fromLocation *Point) ([]Move, error) {
    args := m.Called(fromLocation)
    return args.Get(0).([]Move), args.Error(1)
}

func (m *MockBoard) CalculateMoves(color string) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockBoard) Size() *Point {
    args := m.Called()
    return args.Get(0).(*Point)
}

func (m *MockBoard) Print() string {
    args := m.Called()
    return args.String(0)
}

func (m *MockBoard) State() ([][]*SquareData, bool, bool, bool) {
    args := m.Called()
    return args.Get(0).([][]*SquareData), args.Bool(1), args.Bool(2), args.Bool(3)
}

func (m *MockBoard) pointOutOfBounds(p *Point) bool {
    args := m.Called(p)
    return args.Bool(0)
}

func (m *MockBoard) pointOnPromotionSquare(p *Point) bool {
    args := m.Called(p)
    return args.Bool(0)
}

func Test_NewSimpleBoard_DefaultFen(t *testing.T) {
    s, err := createSimpleBoardWithDefaultPieceLocations()
	assert.Nil(t, err)

	for y := 2; y <= 5; y++ {
		for x := 0; x <= 7; x++ {
			piece, err := s.getPiece(&Point{x, y})
			assert.Nil(t, err)
			assert.Nil(t, piece)
		}
	}

	for _, y := range []int{1, 6} {
		for x := 0; x <= 7; x++ {
			piece, err := s.getPiece(&Point{x, y})
			assert.Nil(t, err)
			_, ok := piece.(*Pawn)
			assert.True(t, ok)
		}
	}

	for _, y := range []int{0, 7} {
		piece, err := s.getPiece(&Point{0, y})
		assert.Nil(t, err)
		_, ok := piece.(*Rook)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{1, y})
		assert.Nil(t, err)
		_, ok = piece.(*Knight)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{2, y})
		assert.Nil(t, err)
		_, ok = piece.(*Bishop)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{3, y})
		assert.Nil(t, err)
		_, ok = piece.(*Queen)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{4, y})
		assert.Nil(t, err)
		_, ok = piece.(*King)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{5, y})
		assert.Nil(t, err)
		_, ok = piece.(*Bishop)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{6, y})
		assert.Nil(t, err)
		_, ok = piece.(*Knight)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{7, y})
		assert.Nil(t, err)
		_, ok = piece.(*Rook)
		assert.True(t, ok)
	}
}

func Test_CalculateMoves_default(t *testing.T) {
    s, err := createSimpleBoardWithDefaultPieceLocations()
    assert.Nil(t, err)

    s.CalculateMoves("white")
    _, check, checkmate, stalemate := s.State()

    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 20, moveCount)
    assert.False(t, check)
    assert.False(t, checkmate)
    assert.False(t, stalemate)
}

func Test_CalculateMoves_check(t *testing.T) {
    s, err := newSimpleBoard(&Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(&Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(&Point{0, 1}, newQueen("black"))
    s.setPiece(&Point{0, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves("white")
    _, check, checkmate, stalemate := s.State()

    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 1, moveCount)
    assert.True(t, check)
    assert.False(t, checkmate)
    assert.False(t, stalemate)
}

func Test_CalculateMoves_checkmate(t *testing.T) {
    s, err := newSimpleBoard(&Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(&Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(&Point{0, 1}, newQueen("black"))
    s.setPiece(&Point{0, 2}, newKing("black", false, 0, -1))

    s.CalculateMoves("white")
    _, check, checkmate, stalemate := s.State()

    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 0, moveCount)
    assert.True(t, check)
    assert.True(t, checkmate)
    assert.False(t, stalemate)
}

func Test_CalculateMoves_stalemate(t *testing.T) {
    s, err := newSimpleBoard(&Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(&Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(&Point{1, 2}, newQueen("black"))
    s.setPiece(&Point{0, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves("white")
    _, check, checkmate, stalemate := s.State()

    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 0, moveCount)
    assert.False(t, check)
    assert.False(t, checkmate)
    assert.True(t, stalemate)
}

func Test_CalculateMoves_noCastleThroughCheck(t *testing.T) {
    s, err := newSimpleBoard(&Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(&Point{4, 0}, newKing("white", false, 0, 1))
    s.setPiece(&Point{0, 0}, newRook("white", false))
    s.setPiece(&Point{3, 7}, newRook("black", false))
    s.setPiece(&Point{4, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves("white")
    _, check, checkmate, stalemate := s.State()


    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 13, moveCount)
    assert.False(t, check)
    assert.False(t, checkmate)
    assert.False(t, stalemate)
}

