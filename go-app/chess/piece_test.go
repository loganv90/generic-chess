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
	moveFactoryInstance = moveFactory

	pawn, err := newPawn("white", false, 0, 1)
	assert.Nil(t, err)

	moves := pawn.moves(board, 3, 3)
	assert.Len(t, moves, 2)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func TestPawnMovesWhenMoved(t *testing.T) {
	board := &mockBoard{}
	board.On("getPiece", mock.Anything, mock.Anything).Return(nil, nil)
	board.On("xLen").Return(8)
	board.On("yLen").Return(8)
	board.On("possibleEnPassants", mock.Anything, mock.Anything, mock.Anything).Return([]*enPassant{})

	moveFactory := &mockMoveFactory{}
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 4).Return(nil, nil)
	moveFactoryInstance = moveFactory

	pawn, err := newPawn("white", true, 0, 1)
	assert.Nil(t, err)

	moves := pawn.moves(board, 3, 3)
	assert.Len(t, moves, 1)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func TestPawnMovesWhenCapturing(t *testing.T) {
	board := &mockBoard{}
	board.On("getPiece", 2, 4).Return(nil, nil)
	board.On("getPiece", 3, 4).Return(nil, nil)
	board.On("getPiece", 3, 5).Return(nil, nil)
	board.On("xLen").Return(8)
	board.On("yLen").Return(8)
	board.On("possibleEnPassants", mock.Anything, mock.Anything, mock.Anything).Return([]*enPassant{})

	moveFactory := &mockMoveFactory{}
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 5).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 4, 4).Return(nil, nil)
	moveFactoryInstance = moveFactory

	pawn, err := newPawn("white", false, 0, 1)
	assert.Nil(t, err)

	blackPawn, err := newPawn("black", false, 0, 1)
	assert.Nil(t, err)
	board.On("getPiece", 4, 4).Return(blackPawn, nil)

	moves := pawn.moves(board, 3, 3)
	assert.Len(t, moves, 3)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func TestPawnMovesWhenCapturingEnPassant(t *testing.T) {
	board := &mockBoard{}
	board.On("getPiece", 2, 4).Return(nil, nil)
	board.On("getPiece", 3, 4).Return(nil, nil)
	board.On("getPiece", 3, 5).Return(nil, nil)
	board.On("getPiece", 4, 4).Return(nil, nil)
	board.On("xLen").Return(8)
	board.On("yLen").Return(8)
	board.On("possibleEnPassants", "white", 4, 4).Return([]*enPassant{{4, 4, 3, 4}})
	board.On("possibleEnPassants", "white", 2, 4).Return([]*enPassant{})

	moveFactory := &mockMoveFactory{}
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 5).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 4, 4).Return(nil, nil)
	moveFactoryInstance = moveFactory

	pawn, err := newPawn("white", false, 0, 1)
	assert.Nil(t, err)

	moves := pawn.moves(board, 3, 3)
	assert.Len(t, moves, 3)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func TestKnightMoves(t *testing.T) {
	board := &mockBoard{}
	board.On("getPiece", 2, 3).Return(nil, nil)
	board.On("getPiece", 3, 2).Return(nil, nil)
	board.On("getPiece", 0, 3).Return(nil, nil)
	board.On("getPiece", 3, 0).Return(nil, nil)
	board.On("xLen").Return(8)
	board.On("yLen").Return(8)

	moveFactory := &mockMoveFactory{}
	moveFactory.On("newSimpleMove", board, 1, 1, 2, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 3, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 0, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 3, 0).Return(nil, nil)
	moveFactoryInstance = moveFactory

	knight, err := newKnight("white")
	assert.Nil(t, err)

	moves := knight.moves(board, 1, 1)
	assert.Len(t, moves, 4)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func TestBishopMoves(t *testing.T) {
	board := &mockBoard{}
	board.On("getPiece", 2, 2).Return(nil, nil)
	board.On("getPiece", 3, 3).Return(nil, nil)
	board.On("getPiece", 4, 4).Return(nil, nil)
	board.On("getPiece", 5, 5).Return(nil, nil)
	board.On("getPiece", 6, 6).Return(nil, nil)
	board.On("getPiece", 7, 7).Return(nil, nil)
	board.On("getPiece", 0, 2).Return(nil, nil)
	board.On("getPiece", 2, 0).Return(nil, nil)
	board.On("getPiece", 0, 0).Return(nil, nil)
	board.On("xLen").Return(8)
	board.On("yLen").Return(8)

	moveFactory := &mockMoveFactory{}
	moveFactory.On("newSimpleMove", board, 1, 1, 2, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 3, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 4, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 5, 5).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 6, 6).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 7, 7).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 0, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 2, 0).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 0, 0).Return(nil, nil)
	moveFactoryInstance = moveFactory

	bishop, err := newBishop("white")
	assert.Nil(t, err)

	moves := bishop.moves(board, 1, 1)
	assert.Len(t, moves, 9)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func TestRookMoves(t *testing.T) {
	board := &mockBoard{}
	board.On("getPiece", 2, 1).Return(nil, nil)
	board.On("getPiece", 3, 1).Return(nil, nil)
	board.On("getPiece", 4, 1).Return(nil, nil)
	board.On("getPiece", 5, 1).Return(nil, nil)
	board.On("getPiece", 6, 1).Return(nil, nil)
	board.On("getPiece", 7, 1).Return(nil, nil)
	board.On("getPiece", 1, 2).Return(nil, nil)
	board.On("getPiece", 1, 3).Return(nil, nil)
	board.On("getPiece", 1, 4).Return(nil, nil)
	board.On("getPiece", 1, 5).Return(nil, nil)
	board.On("getPiece", 1, 6).Return(nil, nil)
	board.On("getPiece", 1, 7).Return(nil, nil)
	board.On("getPiece", 0, 1).Return(nil, nil)
	board.On("getPiece", 1, 0).Return(nil, nil)
	board.On("xLen").Return(8)
	board.On("yLen").Return(8)

	moveFactory := &mockMoveFactory{}
	moveFactory.On("newSimpleMove", board, 1, 1, 2, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 3, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 4, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 5, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 6, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 7, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 5).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 6).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 7).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 0, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 0).Return(nil, nil)
	moveFactoryInstance = moveFactory

	rook, err := newRook("white", false)
	assert.Nil(t, err)

	moves := rook.moves(board, 1, 1)
	assert.Len(t, moves, 14)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func TestQueenMoves(t *testing.T) {
	board := &mockBoard{}
	board.On("getPiece", 2, 2).Return(nil, nil)
	board.On("getPiece", 3, 3).Return(nil, nil)
	board.On("getPiece", 4, 4).Return(nil, nil)
	board.On("getPiece", 5, 5).Return(nil, nil)
	board.On("getPiece", 6, 6).Return(nil, nil)
	board.On("getPiece", 7, 7).Return(nil, nil)
	board.On("getPiece", 0, 2).Return(nil, nil)
	board.On("getPiece", 2, 0).Return(nil, nil)
	board.On("getPiece", 0, 0).Return(nil, nil)
	board.On("getPiece", 2, 1).Return(nil, nil)
	board.On("getPiece", 3, 1).Return(nil, nil)
	board.On("getPiece", 4, 1).Return(nil, nil)
	board.On("getPiece", 5, 1).Return(nil, nil)
	board.On("getPiece", 6, 1).Return(nil, nil)
	board.On("getPiece", 7, 1).Return(nil, nil)
	board.On("getPiece", 1, 2).Return(nil, nil)
	board.On("getPiece", 1, 3).Return(nil, nil)
	board.On("getPiece", 1, 4).Return(nil, nil)
	board.On("getPiece", 1, 5).Return(nil, nil)
	board.On("getPiece", 1, 6).Return(nil, nil)
	board.On("getPiece", 1, 7).Return(nil, nil)
	board.On("getPiece", 0, 1).Return(nil, nil)
	board.On("getPiece", 1, 0).Return(nil, nil)
	board.On("xLen").Return(8)
	board.On("yLen").Return(8)

	moveFactory := &mockMoveFactory{}
	moveFactory.On("newSimpleMove", board, 1, 1, 2, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 3, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 4, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 5, 5).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 6, 6).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 7, 7).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 0, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 2, 0).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 0, 0).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 2, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 3, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 4, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 5, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 6, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 7, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 5).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 6).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 7).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 0, 1).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 1, 1, 1, 0).Return(nil, nil)
	moveFactoryInstance = moveFactory

	queen, err := newQueen("white")
	assert.Nil(t, err)

	moves := queen.moves(board, 1, 1)
	assert.Len(t, moves, 23)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func TestKingMovesWhenCanCastleAndUnmoved(t *testing.T) {
	board := &mockBoard{}
	board.On("getPiece", 2, 2).Return(nil, nil)
	board.On("getPiece", 3, 2).Return(nil, nil)
	board.On("getPiece", 4, 2).Return(nil, nil)
	board.On("getPiece", 2, 3).Return(nil, nil)
	board.On("getPiece", 4, 3).Return(nil, nil)
	board.On("getPiece", 2, 4).Return(nil, nil)
	board.On("getPiece", 3, 4).Return(nil, nil)
	board.On("getPiece", 4, 4).Return(nil, nil)
	board.On("getPiece", 5, 3).Return(nil, nil)
	board.On("getPiece", 6, 3).Return(nil, nil)
	board.On("getPiece", 1, 3).Return(nil, nil)
	board.On("xLen").Return(8)
	board.On("yLen").Return(8)

	moveFactory := &mockMoveFactory{}
	moveFactory.On("newSimpleMove", board, 3, 3, 2, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 4, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 2, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 4, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 2, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 4, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 0, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 7, 3).Return(nil, nil)
	moveFactoryInstance = moveFactory

	king, err := newKing("white", false, 0, 1)
	assert.Nil(t, err)

	rook, err := newRook("white", false)
	assert.Nil(t, err)
	board.On("getPiece", 0, 3).Return(rook, nil)
	board.On("getPiece", 7, 3).Return(rook, nil)

	moves := king.moves(board, 3, 3)
	assert.Len(t, moves, 10)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func TestKingMovesWhenCanCastleAndMoved(t *testing.T) {
	board := &mockBoard{}
	board.On("getPiece", 2, 2).Return(nil, nil)
	board.On("getPiece", 3, 2).Return(nil, nil)
	board.On("getPiece", 4, 2).Return(nil, nil)
	board.On("getPiece", 2, 3).Return(nil, nil)
	board.On("getPiece", 4, 3).Return(nil, nil)
	board.On("getPiece", 2, 4).Return(nil, nil)
	board.On("getPiece", 3, 4).Return(nil, nil)
	board.On("getPiece", 4, 4).Return(nil, nil)
	board.On("xLen").Return(8)
	board.On("yLen").Return(8)

	moveFactory := &mockMoveFactory{}
	moveFactory.On("newSimpleMove", board, 3, 3, 2, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 4, 2).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 2, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 4, 3).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 2, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 3, 4).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, 3, 3, 4, 4).Return(nil, nil)
	moveFactoryInstance = moveFactory

	king, err := newKing("white", true, 0, 1)
	assert.Nil(t, err)

	moves := king.moves(board, 3, 3)
	assert.Len(t, moves, 8)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}
