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
	return args.Get(0).(piece), args.Error(1)
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

func TestSimpleMove(t *testing.T) {
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

	simpleMove, err := newSimpleMove(board, 0, 0, 1, 1)
	assert.Nil(t, err)

	err = simpleMove.execute()
	assert.Nil(t, err)

	err = simpleMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
}
