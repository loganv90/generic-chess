package chess

import (
    "testing"

    "github.com/stretchr/testify/assert"
    //"github.com/stretchr/testify/mock"
)

func Test_SimpleMove2(t *testing.T) {
    white := 0
    from := Point{0, 0}
    to := Point{1, 1}
    start := Point{2, 2}
    end := Point{3, 3}
    target := Point{4, 4}
    risk := Point{5, 5}
    board := &MockBoard{}
    piece := &MockPiece{}
    newPiece := &MockPiece{}
    capturedPiece := &MockPiece{}

	board.On("getPiece", from).Return(piece, true)
	board.On("getPiece", to).Return(capturedPiece, true)
	board.On("getEnPassant2", white).Return(target, risk, nil)
    board.On("getVulnerables2", white).Return(start, end, nil)
	piece.On("copy").Return(newPiece, nil)
    simpleMove, err := createSimpleMove(board, white, from, to, nil)
    assert.Nil(t, err)

	board.On("setPiece", from, nil).Return(true)
	board.On("setPiece", to, newPiece).Return(true)
	board.On("setEnPassant2", white, Point{-1, -1}, Point{-1, -1}).Return(nil)
    board.On("setVulnerables2", white, Point{-1, -1}, Point{-1, -1}).Return(nil)
    err = simpleMove.execute()
    assert.Nil(t, err)

	board.On("setPiece", from, piece).Return(true)
	board.On("setPiece", to, capturedPiece).Return(true)
	board.On("setEnPassant2", white, target, risk).Return(nil)
    board.On("setVulnerables2", white, start, end).Return(nil)
    err = simpleMove.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
    newPiece.AssertExpectations(t)
}

// TODO add the rest of the tests

