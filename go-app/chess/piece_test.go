package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPiece struct {
	mock.Mock
}

func (m *mockPiece) getColor() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockPiece) movedCopy() piece {
	args := m.Called()
	return args.Get(0).(piece)
}

func (m *mockPiece) moves(board board, x int, y int) []move {
	args := m.Called(board, x, y)
	return args.Get(0).([]move)
}

func TestPawnMovesWhenUnmoved(t *testing.T) {
	board := &mockBoard{}
	board.On("getPiece", mock.Anything, mock.Anything).Return(nil, nil)
	board.On("xLen").Return(8)
	board.On("yLen").Return(8)
	board.On("possibleEnPassants", mock.Anything, mock.Anything, mock.Anything).Return([]*enPassant{})

	moveFactory := &mockMoveFactory{}
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 5).Return(nil, nil)

	pawn, err := newPawn(moveFactory, "white", false, 0, 1)
	assert.Nil(t, err)

	moves := pawn.moves(board, 3, 3)
	assert.Len(t, moves, 2)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}
