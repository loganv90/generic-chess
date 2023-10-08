package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSimpleBoardWithDefaultFen(t *testing.T) {
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
