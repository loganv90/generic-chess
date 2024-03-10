package chess

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_MoveSimple(t *testing.T) {
    white := 0
    from := Point{0, 0}
    to := Point{1, 1}
    vulnerable := Vulnerable{Point{2, 2}, Point{3, 3}}
    enPassant := EnPassant{Point{4, 4}, Point{5, 5}}
    board := &MockBoard{}
    piece := &MockPiece{}
    newPiece := &MockPiece{}
    capturedPiece := &MockPiece{}

	board.On("getPiece", from).Return(piece, true)
	board.On("getPiece", to).Return(capturedPiece, true)
	board.On("getEnPassant", white).Return(enPassant, nil)
    board.On("getVulnerable", white).Return(vulnerable, nil)
	piece.On("copy").Return(newPiece, nil)
    simpleMove, err := createMoveSimple(board, white, from, to, nil)
    assert.Nil(t, err)
    assert.Equal(t, simpleMove.allyDefense, false)

	board.On("setPiece", from, nil).Return(true)
	board.On("setPiece", to, newPiece).Return(true)
	board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    err = simpleMove.execute()
    assert.Nil(t, err)

	board.On("setPiece", from, piece).Return(true)
	board.On("setPiece", to, capturedPiece).Return(true)
	board.On("setEnPassant", white, enPassant).Return(nil)
    board.On("setVulnerable", white, vulnerable).Return(nil)
    err = simpleMove.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
}

func Test_MovePromotion(t *testing.T) {
    white := 0
    from := Point{0, 0}
    to := Point{1, 1}
    vulnerable := Vulnerable{Point{2, 2}, Point{3, 3}}
    enPassant := EnPassant{Point{4, 4}, Point{5, 5}}
    board := &MockBoard{}
    piece := &MockPiece{}
    newPiece := &MockPiece{}
    capturedPiece := &MockPiece{}

	board.On("getPiece", from).Return(piece, true)
	board.On("getPiece", to).Return(capturedPiece, true)
	board.On("getEnPassant", white).Return(enPassant, nil)
    board.On("getVulnerable", white).Return(vulnerable, nil)
    promotionMove, err := createMoveSimple(board, white, from, to, newPiece)
    assert.Nil(t, err)
    assert.Equal(t, promotionMove.allyDefense, false)

	board.On("setPiece", from, nil).Return(true)
	board.On("setPiece", to, newPiece).Return(true)
	board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    err = promotionMove.execute()
    assert.Nil(t, err)

	board.On("setPiece", from, piece).Return(true)
	board.On("setPiece", to, capturedPiece).Return(true)
	board.On("setEnPassant", white, enPassant).Return(nil)
    board.On("setVulnerable", white, vulnerable).Return(nil)
    err = promotionMove.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
}

func Test_MoveRevealEnPassant(t *testing.T) {
    white := 0
    from := Point{0, 0}
    to := Point{1, 1}
    vulnerable := Vulnerable{Point{2, 2}, Point{3, 3}}
    enPassant := EnPassant{Point{4, 4}, Point{5, 5}}
    board := &MockBoard{}
    piece := &MockPiece{}
    newPiece := &MockPiece{}
    capturedPiece := &MockPiece{}

    newEnPassant := EnPassant{Point{6, 6}, to}

	board.On("getPiece", from).Return(piece, true)
	board.On("getPiece", to).Return(capturedPiece, true)
	board.On("getEnPassant", white).Return(enPassant, nil)
    board.On("getVulnerable", white).Return(vulnerable, nil)
	piece.On("copy").Return(newPiece, nil)
    revealEnPassantMove, err := createMoveRevealEnPassant(board, white, from, to, nil, newEnPassant)
    assert.Nil(t, err)
    assert.Equal(t, revealEnPassantMove.allyDefense, false)

	board.On("setPiece", from, nil).Return(true)
	board.On("setPiece", to, newPiece).Return(true)
	board.On("setEnPassant", white, newEnPassant).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    err = revealEnPassantMove.execute()
    assert.Nil(t, err)

	board.On("setPiece", from, piece).Return(true)
	board.On("setPiece", to, capturedPiece).Return(true)
	board.On("setEnPassant", white, enPassant).Return(nil)
    board.On("setVulnerable", white, vulnerable).Return(nil)
    err = revealEnPassantMove.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
}

func Test_MoveCaptureEnPassant(t *testing.T) {
    white := 0
    from := Point{0, 0}
    to := Point{1, 1}
    vulnerable := Vulnerable{Point{2, 2}, Point{3, 3}}
    enPassant := EnPassant{Point{4, 4}, Point{5, 5}}
    board := &MockBoard{}
    piece := &MockPiece{}
    newPiece := &MockPiece{}
    capturedPiece := &MockPiece{}

    enPassantRisk1 := Point{7, 7}
    enPassantRisk2 := Point{9, 9}
    enPassant1 := EnPassant{Point{6, 6}, enPassantRisk1}
    enPassant2 := EnPassant{Point{8, 8}, enPassantRisk2}
    enPassantCapturedPiece1 := &MockPiece{}
    enPassantCapturedPiece2 := &MockPiece{}

	board.On("getPiece", from).Return(piece, true)
	board.On("getPiece", to).Return(capturedPiece, true)
	board.On("getEnPassant", white).Return(enPassant, nil)
    board.On("getVulnerable", white).Return(vulnerable, nil)
	piece.On("copy").Return(newPiece, nil)
    board.On("possibleEnPassant", white, to).Return([]EnPassant{enPassant1, enPassant2}, nil)
    board.On("getPiece", enPassantRisk1).Return(enPassantCapturedPiece1, true)
    board.On("getPiece", enPassantRisk2).Return(enPassantCapturedPiece2, true)
    captureEnPassantMove, err := createMoveCaptureEnPassant(board, white, from, to, nil)
    assert.Nil(t, err)
    assert.Equal(t, captureEnPassantMove.allyDefense, false)

	board.On("setPiece", from, nil).Return(true)
	board.On("setPiece", to, newPiece).Return(true)
	board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
	board.On("setPiece", enPassantRisk1, nil).Return(true)
	board.On("setPiece", enPassantRisk2, nil).Return(true)
    err = captureEnPassantMove.execute()
    assert.Nil(t, err)

	board.On("setPiece", from, piece).Return(true)
	board.On("setPiece", to, capturedPiece).Return(true)
	board.On("setEnPassant", white, enPassant).Return(nil)
    board.On("setVulnerable", white, vulnerable).Return(nil)
    board.On("setPiece", enPassantRisk1, enPassantCapturedPiece1).Return(true)
    board.On("setPiece", enPassantRisk2, enPassantCapturedPiece2).Return(true)
    err = captureEnPassantMove.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
}

func Test_MoveAllyDefense(t *testing.T) {
    white := 0
    from := Point{0, 0}
    to := Point{1, 1}
    board := &MockBoard{}

    allyDefenseMove, err := createMoveAllyDefense(board, white, from, to)
    assert.Nil(t, err)
    assert.Equal(t, allyDefenseMove.allyDefense, true)

    board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    err = allyDefenseMove.execute()
    assert.Nil(t, err)

    board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    err = allyDefenseMove.undo()
    assert.Nil(t, err)

    board.AssertExpectations(t)
}

func Test_MoveCastle(t *testing.T) {
    white := 0
    from := Point{0, 0}
    to := Point{1, 1}
    vulnerable := Vulnerable{Point{2, 2}, Point{3, 3}}
    enPassant := EnPassant{Point{4, 4}, Point{5, 5}}
    toKing := Point{6, 6}
    toRook := Point{7, 7}
    newVulnerable := Vulnerable{Point{8, 8}, Point{9, 9}}
    board := &MockBoard{}
    king := &MockPiece{}
    newKing := &MockPiece{}
    rook := &MockPiece{}
    newRook := &MockPiece{}

	board.On("getPiece", from).Return(king, true)
	board.On("getPiece", to).Return(rook, true)
	board.On("getEnPassant", white).Return(enPassant, nil)
    board.On("getVulnerable", white).Return(vulnerable, nil)
	king.On("copy").Return(newKing, nil)
	rook.On("copy").Return(newRook, nil)
    castleMove, err := createMoveCastle(board, white, from, to, toKing, toRook, newVulnerable)
    assert.Nil(t, err)
    assert.Equal(t, castleMove.allyDefense, false)

	board.On("setPiece", from, nil).Return(true)
	board.On("setPiece", to, nil).Return(true)
	board.On("setPiece", toKing, newKing).Return(true)
	board.On("setPiece", toRook, newRook).Return(true)
	board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, newVulnerable).Return(nil)
    err = castleMove.execute()
    assert.Nil(t, err)

	board.On("setPiece", toKing, nil).Return(true)
	board.On("setPiece", toRook, nil).Return(true)
	board.On("setPiece", from, king).Return(true)
	board.On("setPiece", to, rook).Return(true)
	board.On("setEnPassant", white, enPassant).Return(nil)
    board.On("setVulnerable", white, vulnerable).Return(nil)
    err = castleMove.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
    king.AssertExpectations(t)
    rook.AssertExpectations(t)
}

