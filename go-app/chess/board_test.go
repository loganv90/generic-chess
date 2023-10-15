package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBoard struct {
	mock.Mock
}

func (m *mockBoard) getPiece(x int, y int) (piece, error) {
	args := m.Called(x, y)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(piece), args.Error(1)
	}
}

func (m *mockBoard) setPiece(x int, y int, piece piece) error {
	args := m.Called(x, y, piece)
	return args.Error(0)
}

func (m *mockBoard) getEnPassant(color string) (*enPassant, error) {
	args := m.Called(color)
	return args.Get(0).(*enPassant), args.Error(1)
}

func (m *mockBoard) setEnPassant(color string, enPassant *enPassant) {
	m.Called(color, enPassant)
}

func (m *mockBoard) clrEnPassant(color string) {
	m.Called(color)
}

func (m *mockBoard) possibleEnPassants(color string, xTarget int, yTarget int) []*enPassant {
	args := m.Called(color, xTarget, yTarget)
	return args.Get(0).([]*enPassant)
}

func (m *mockBoard) moves(x int, y int) []move {
	args := m.Called(x, y)
	return args.Get(0).([]move)
}

func (m *mockBoard) increment() {
	m.Called()
}

func (m *mockBoard) decrement() {
	m.Called()
}

func (m *mockBoard) xLen() int {
	args := m.Called()
	return args.Int(0)
}

func (m *mockBoard) yLen() int {
	args := m.Called()
	return args.Int(0)
}

func Test_NewSimpleBoard_DefaultFen(t *testing.T) {
	s, err := newSimpleBoard(
		[]string{"white", "black"},
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	)
	assert.Nil(t, err)

	for y := 2; y <= 5; y++ {
		for x := 0; x <= 7; x++ {
			piece, err := s.getPiece(x, y)
			assert.Nil(t, err)
			assert.Nil(t, piece)
		}
	}

	for _, y := range []int{1, 6} {
		for x := 0; x <= 7; x++ {
			piece, err := s.getPiece(x, y)
			assert.Nil(t, err)
			_, ok := piece.(*pawn)
			assert.True(t, ok)
		}
	}

	for _, y := range []int{0, 7} {
		piece, err := s.getPiece(0, y)
		assert.Nil(t, err)
		_, ok := piece.(*rook)
		assert.True(t, ok)

		piece, err = s.getPiece(1, y)
		assert.Nil(t, err)
		_, ok = piece.(*knight)
		assert.True(t, ok)

		piece, err = s.getPiece(2, y)
		assert.Nil(t, err)
		_, ok = piece.(*bishop)
		assert.True(t, ok)

		piece, err = s.getPiece(3, y)
		assert.Nil(t, err)
		_, ok = piece.(*queen)
		assert.True(t, ok)

		piece, err = s.getPiece(4, y)
		assert.Nil(t, err)
		_, ok = piece.(*king)
		assert.True(t, ok)

		piece, err = s.getPiece(5, y)
		assert.Nil(t, err)
		_, ok = piece.(*bishop)
		assert.True(t, ok)

		piece, err = s.getPiece(6, y)
		assert.Nil(t, err)
		_, ok = piece.(*knight)
		assert.True(t, ok)

		piece, err = s.getPiece(7, y)
		assert.Nil(t, err)
		_, ok = piece.(*rook)
		assert.True(t, ok)
	}
}
