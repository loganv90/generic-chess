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

func (m *mockMoveFactory) newRevealEnPassantMove(
	b board,
	xFrom int,
	yFrom int,
	xTo int,
	yTo int,
	xTarget int,
	yTarget int,
) (*revealEnPassantMove, error) {
	args := m.Called(b, xFrom, yFrom, xTo, yTo, xTarget, yTarget)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*revealEnPassantMove), args.Error(1)
	}
}

func (m *mockMoveFactory) newCaptureEnPassantMove(
	b board,
	xFrom int,
	yFrom int,
	xTo int,
	yTo int,
) (*captureEnPassantMove, error) {
	args := m.Called(b, xFrom, yFrom, xTo, yTo)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*captureEnPassantMove), args.Error(1)
	}
}

func (m *mockMoveFactory) newCastleMove(
	b board,
	xFrom int,
	yFrom int,
	xTo int,
	yTo int,
	xToKing int,
	yToKing int,
	xToRook int,
	yToRook int,
) (*castleMove, error) {
	args := m.Called(b, xFrom, yFrom, xTo, yTo, xToKing, yToKing, xToRook, yToRook)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*castleMove), args.Error(1)
	}
}

type mockMove struct {
	mock.Mock
}

func (m *mockMove) execute() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockMove) undo() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockMove) getAction() *action {
	args := m.Called()
	return args.Get(0).(*action)
}

func Test_SimpleMove(t *testing.T) {
	board := &mockBoard{}
	piece := &mockPiece{}
	newPiece := &mockPiece{}
	capturedPiece := &mockPiece{}
	en := &enPassant{}

	board.On("getPiece", 0, 0).Return(piece, nil)
	piece.On("movedCopy").Return(newPiece)
	board.On("getPiece", 1, 1).Return(capturedPiece, nil)
	board.On("getEnPassant", "white").Return(en, nil)
	piece.On("getColor").Return("white")
	simpleMove, err := moveFactoryInstance.newSimpleMove(board, 0, 0, 1, 1)
	assert.Nil(t, err)

	board.On("setPiece", 0, 0, nil).Return(nil)
	board.On("setPiece", 1, 1, newPiece).Return(nil)
	board.On("clrEnPassant", "white").Return()
	board.On("increment").Return()
	err = simpleMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", 0, 0, piece).Return(nil)
	board.On("setPiece", 1, 1, capturedPiece).Return(nil)
	board.On("setEnPassant", "white", en).Return()
	board.On("decrement").Return()
	err = simpleMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
}

func Test_RevealEnPassantMove(t *testing.T) {
	board := &mockBoard{}
	piece := &mockPiece{}
	newPiece := &mockPiece{}
	capturedPiece := &mockPiece{}
	en := &enPassant{}

	board.On("getPiece", 0, 0).Return(piece, nil)
	piece.On("movedCopy").Return(newPiece)
	board.On("getPiece", 2, 2).Return(capturedPiece, nil)
	board.On("getEnPassant", "white").Return(en, nil)
	piece.On("getColor").Return("white")
	revealEnPassantMove, err := moveFactoryInstance.newRevealEnPassantMove(board, 0, 0, 2, 2, 1, 1)
	assert.Nil(t, err)

	board.On("setPiece", 0, 0, nil).Return(nil)
	board.On("setPiece", 2, 2, newPiece).Return(nil)
	board.On("setEnPassant", "white", &enPassant{1, 1, 2, 2}).Return()
	board.On("increment").Return()
	err = revealEnPassantMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", 0, 0, piece).Return(nil)
	board.On("setPiece", 2, 2, capturedPiece).Return(nil)
	board.On("setEnPassant", "white", en).Return()
	board.On("decrement").Return()
	err = revealEnPassantMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
}

func Test_CaptureEnPassantMove(t *testing.T) {
	board := &mockBoard{}
	piece := &mockPiece{}
	newPiece := &mockPiece{}
	capturedPiece := &mockPiece{}
	en := &enPassant{}
	encPiece := &mockPiece{}

	board.On("getPiece", 0, 0).Return(piece, nil)
	piece.On("movedCopy").Return(newPiece)
	board.On("getPiece", 1, 1).Return(capturedPiece, nil)
	board.On("getEnPassant", "white").Return(en, nil)
	piece.On("getColor").Return("white")
	board.On("possibleEnPassants", "white", 1, 1).Return([]*enPassant{{1, 1, 2, 2}})
	board.On("getPiece", 2, 2).Return(encPiece, nil)
	captureEnPassantMove, err := moveFactoryInstance.newCaptureEnPassantMove(board, 0, 0, 1, 1)
	assert.Nil(t, err)

	assert.NotNil(t, captureEnPassantMove)

	board.On("setPiece", 0, 0, nil).Return(nil)
	board.On("setPiece", 1, 1, newPiece).Return(nil)
	board.On("setPiece", 2, 2, nil).Return(nil)
	board.On("clrEnPassant", "white").Return()
	board.On("increment").Return()
	err = captureEnPassantMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", 0, 0, piece).Return(nil)
	board.On("setPiece", 1, 1, capturedPiece).Return(nil)
	board.On("setPiece", 2, 2, encPiece).Return(nil)
	board.On("setEnPassant", "white", en).Return()
	board.On("decrement").Return()
	err = captureEnPassantMove.undo()
	assert.Nil(t, err)
}

func Test_CastleMove(t *testing.T) {
	board := &mockBoard{}
	king := &mockPiece{}
	newKing := &mockPiece{}
	rook := &mockPiece{}
	newRook := &mockPiece{}
	en := &enPassant{}

	board.On("getPiece", 0, 0).Return(king, nil)
	board.On("getPiece", 1, 1).Return(rook, nil)
	king.On("movedCopy").Return(newKing)
	rook.On("movedCopy").Return(newRook)
	board.On("getEnPassant", "white").Return(en, nil)
	king.On("getColor").Return("white")
	castleMove, err := moveFactoryInstance.newCastleMove(board, 0, 0, 1, 1, 2, 2, 3, 3)
	assert.Nil(t, err)

	board.On("setPiece", 0, 0, nil).Return(nil)
	board.On("setPiece", 1, 1, nil).Return(nil)
	board.On("setPiece", 2, 2, newKing).Return(nil)
	board.On("setPiece", 3, 3, newRook).Return(nil)
	board.On("clrEnPassant", "white").Return()
	board.On("increment").Return()
	err = castleMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", 0, 0, king).Return(nil)
	board.On("setPiece", 1, 1, rook).Return(nil)
	board.On("setPiece", 2, 2, nil).Return(nil)
	board.On("setPiece", 3, 3, nil).Return(nil)
	board.On("setEnPassant", "white", en).Return()
	board.On("decrement").Return()
	err = castleMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	king.AssertExpectations(t)
	rook.AssertExpectations(t)
}
