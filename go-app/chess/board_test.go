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

func (m *MockBoard) getEnPassant(color string) (*EnPassant, error) {
	args := m.Called(color)
	return args.Get(0).(*EnPassant), args.Error(1)
}

func (m *MockBoard) setEnPassant(color string, enPassant *EnPassant) {
	m.Called(color, enPassant)
}

func (m *MockBoard) clrEnPassant(color string) {
	m.Called(color)
}

func (m *MockBoard) possibleEnPassants(color string, target *Point) []*EnPassant {
	args := m.Called(color, target)
	return args.Get(0).([]*EnPassant)
}

func (m *MockBoard) moves(location *Point) []Move {
	args := m.Called(location)
	return args.Get(0).([]Move)
}

func (m *MockBoard) increment() {
	m.Called()
}

func (m *MockBoard) decrement() {
	m.Called()
}

func (m *MockBoard) xLen() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockBoard) yLen() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockBoard) print() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockBoard) turn() string {
    args := m.Called()
    return args.String(0)
}

func (m *MockBoard) squares() [][]*SquareData {
    args := m.Called()
    return args.Get(0).([][]*SquareData)
}

func (m *MockBoard) checkmate(color string) bool {
    args := m.Called(color)
    return args.Bool(0)
}

func (m *MockBoard) check(color string) bool {
    args := m.Called(color)
    return args.Bool(0)
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
	s, err := newSimpleBoard(
		[]string{"white", "black"},
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	)
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
