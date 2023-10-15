package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockMoveFactory struct {
	mock.Mock
}

func (m *mockMoveFactory) newSimpleMove(
	b board,
	xFrom int,
	yFrom int,
	xTo int,
	yTo int,
) (*simpleMove, error) {
	args := m.Called(b, xFrom, yFrom, xTo, yTo)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*simpleMove), args.Error(1)
	}
}

func Test_SimpleMove(t *testing.T) {
	board := &mockBoard{}
	piece := &mockPiece{}
	newPiece := &mockPiece{}
	capturedPiece := &mockPiece{}
	enPassant := &enPassant{}

	board.On("getPiece", 0, 0).Return(piece, nil)
	board.On("getPiece", 1, 1).Return(capturedPiece, nil)
	board.On("getEnPassant", "white").Return(enPassant, nil)
	piece.On("movedCopy").Return(newPiece)
	piece.On("getColor").Return("white")
	board.On("setPiece", 0, 0, nil).Return(nil)
	board.On("setPiece", 1, 1, newPiece).Return(nil)
	board.On("setPiece", 0, 0, piece).Return(nil)
	board.On("setPiece", 1, 1, capturedPiece).Return(nil)
	board.On("clrEnPassant", "white").Return()
	board.On("setEnPassant", "white", enPassant).Return()
	board.On("increment").Return()
	board.On("decrement").Return()

	simpleMove, err := moveFactoryInstance.newSimpleMove(board, 0, 0, 1, 1)
	assert.Nil(t, err)

	err = simpleMove.execute()
	assert.Nil(t, err)

	err = simpleMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
}
