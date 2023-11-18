package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMoveFactory struct {
	mock.Mock
}

func (m *MockMoveFactory) newSimpleMove(b Board, fromLocation *Point, toLocation *Point) (*SimpleMove, error) {
	args := m.Called(b, fromLocation, toLocation)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*SimpleMove), args.Error(1)
	}
}

func (m *MockMoveFactory) newRevealEnPassantMove(b Board, fromLocation *Point, toLocation *Point, target *Point) (*RevealEnPassantMove, error) {
    args := m.Called(b, fromLocation, toLocation, target)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*RevealEnPassantMove), args.Error(1)
	}
}

func (m *MockMoveFactory) newCaptureEnPassantMove(b Board, fromLocation *Point, toLocation *Point) (*CaptureEnPassantMove, error) {
    args := m.Called(b, fromLocation, toLocation)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*CaptureEnPassantMove), args.Error(1)
	}
}

func (m *MockMoveFactory) newCastleMove(b Board, fromLocation *Point, toLocation *Point, toKingLocation *Point, toRookLocation *Point) (*CastleMove, error) {
    args := m.Called(b, fromLocation, toLocation, toKingLocation, toRookLocation)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*CastleMove), args.Error(1)
	}
}

type MockMove struct {
	mock.Mock
}

func (m *MockMove) execute() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMove) undo() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMove) getAction() *Action {
	args := m.Called()
	return args.Get(0).(*Action)
}

func Test_SimpleMove(t *testing.T) {
	board := &MockBoard{}
	piece := &MockPiece{}
	newPiece := &MockPiece{}
	capturedPiece := &MockPiece{}
	en := &EnPassant{}

	board.On("getPiece", &Point{0, 0}).Return(piece, nil)
	piece.On("movedCopy").Return(newPiece)
	board.On("getPiece", &Point{1, 1}).Return(capturedPiece, nil)
	board.On("getEnPassant", "white").Return(en, nil)
	piece.On("getColor").Return("white")
	simpleMove, err := moveFactoryInstance.newSimpleMove(board, &Point{0, 0}, &Point{1, 1})
	assert.Nil(t, err)

	board.On("setPiece", &Point{0, 0}, nil).Return(nil)
	board.On("setPiece", &Point{1, 1}, newPiece).Return(nil)
	board.On("clrEnPassant", "white").Return()
	board.On("increment").Return()
	err = simpleMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", &Point{0, 0}, piece).Return(nil)
	board.On("setPiece", &Point{1, 1}, capturedPiece).Return(nil)
	board.On("setEnPassant", "white", en).Return()
	board.On("decrement").Return()
	err = simpleMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
}

func Test_RevealEnPassantMove(t *testing.T) {
	board := &MockBoard{}
	piece := &MockPiece{}
	newPiece := &MockPiece{}
	capturedPiece := &MockPiece{}
	en := &EnPassant{}

	board.On("getPiece", &Point{0, 0}).Return(piece, nil)
	piece.On("movedCopy").Return(newPiece)
	board.On("getPiece", &Point{2, 2}).Return(capturedPiece, nil)
	board.On("getEnPassant", "white").Return(en, nil)
	piece.On("getColor").Return("white")
	revealEnPassantMove, err := moveFactoryInstance.newRevealEnPassantMove(board, &Point{0, 0}, &Point{2, 2}, &Point{1, 1})
	assert.Nil(t, err)

	board.On("setPiece", &Point{0, 0}, nil).Return(nil)
	board.On("setPiece", &Point{2, 2}, newPiece).Return(nil)
	board.On("setEnPassant", "white", &EnPassant{&Point{1, 1}, &Point{2, 2}}).Return()
	board.On("increment").Return()
	err = revealEnPassantMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", &Point{0, 0}, piece).Return(nil)
	board.On("setPiece", &Point{2, 2}, capturedPiece).Return(nil)
	board.On("setEnPassant", "white", en).Return()
	board.On("decrement").Return()
	err = revealEnPassantMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
}

func Test_CaptureEnPassantMove(t *testing.T) {
	board := &MockBoard{}
	piece := &MockPiece{}
	newPiece := &MockPiece{}
	capturedPiece := &MockPiece{}
	en := &EnPassant{}
	encPiece := &MockPiece{}

	board.On("getPiece", &Point{0, 0}).Return(piece, nil)
	piece.On("movedCopy").Return(newPiece)
	board.On("getPiece", &Point{1, 1}).Return(capturedPiece, nil)
	board.On("getEnPassant", "white").Return(en, nil)
	piece.On("getColor").Return("white")
	board.On("possibleEnPassants", "white", &Point{1, 1}).Return([]*EnPassant{{&Point{1, 1}, &Point{2, 2}}})
	board.On("getPiece", &Point{2, 2}).Return(encPiece, nil)
	captureEnPassantMove, err := moveFactoryInstance.newCaptureEnPassantMove(board, &Point{0, 0}, &Point{1, 1})
	assert.Nil(t, err)

	assert.NotNil(t, captureEnPassantMove)

	board.On("setPiece", &Point{0, 0}, nil).Return(nil)
	board.On("setPiece", &Point{1, 1}, newPiece).Return(nil)
	board.On("setPiece", &Point{2, 2}, nil).Return(nil)
	board.On("clrEnPassant", "white").Return()
	board.On("increment").Return()
	err = captureEnPassantMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", &Point{0, 0}, piece).Return(nil)
	board.On("setPiece", &Point{1, 1}, capturedPiece).Return(nil)
	board.On("setPiece", &Point{2, 2}, encPiece).Return(nil)
	board.On("setEnPassant", "white", en).Return()
	board.On("decrement").Return()
	err = captureEnPassantMove.undo()
	assert.Nil(t, err)
}

func Test_CastleMove(t *testing.T) {
	board := &MockBoard{}
	king := &MockPiece{}
	newKing := &MockPiece{}
	rook := &MockPiece{}
	newRook := &MockPiece{}
	en := &EnPassant{}

	board.On("getPiece", &Point{0, 0}).Return(king, nil)
	board.On("getPiece", &Point{1, 1}).Return(rook, nil)
	king.On("movedCopy").Return(newKing)
	rook.On("movedCopy").Return(newRook)
	board.On("getEnPassant", "white").Return(en, nil)
	king.On("getColor").Return("white")
	castleMove, err := moveFactoryInstance.newCastleMove(board, &Point{0, 0}, &Point{1, 1}, &Point{2, 2}, &Point{3, 3})
	assert.Nil(t, err)

	board.On("setPiece", &Point{0, 0}, nil).Return(nil)
	board.On("setPiece", &Point{1, 1}, nil).Return(nil)
	board.On("setPiece", &Point{2, 2}, newKing).Return(nil)
	board.On("setPiece", &Point{3, 3}, newRook).Return(nil)
	board.On("clrEnPassant", "white").Return()
	board.On("increment").Return()
	err = castleMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", &Point{0, 0}, king).Return(nil)
	board.On("setPiece", &Point{1, 1}, rook).Return(nil)
	board.On("setPiece", &Point{2, 2}, nil).Return(nil)
	board.On("setPiece", &Point{3, 3}, nil).Return(nil)
	board.On("setEnPassant", "white", en).Return()
	board.On("decrement").Return()
	err = castleMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	king.AssertExpectations(t)
	rook.AssertExpectations(t)
}

