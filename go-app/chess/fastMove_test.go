package chess

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_MoveSimple(t *testing.T) {
    white := 0
    black := 1
    from := Point{0, 0}
    to := Point{1, 1}
    vulnerable := Vulnerable{Point{2, 2}, Point{3, 3}}
    enPassant := EnPassant{Point{4, 4}, Point{5, 5}}
    board := &MockBoard{}
    piece := Piece{white, 1}
    newPiece := Piece{white, 5}
    capturedPiece := Piece{black, 1}

	board.On("getEnPassant", white).Return(enPassant, nil)
    board.On("getVulnerable", white).Return(vulnerable, nil)
    simpleMove := createMoveSimple(board, piece, from, capturedPiece, to, Piece{0, 0})
    assert.Equal(t, simpleMove.allyDefense, false)

	board.On("setPiece", from, Piece{0, 0}).Return(true)
	board.On("setPiece", to, newPiece).Return(true)
	board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    err := simpleMove.execute()
    assert.Nil(t, err)

	board.On("setPiece", from, piece).Return(true)
	board.On("setPiece", to, capturedPiece).Return(true)
	board.On("setEnPassant", white, enPassant).Return(nil)
    board.On("setVulnerable", white, vulnerable).Return(nil)
    err = simpleMove.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
}

func Test_MovePromotion(t *testing.T) {
    white := 0
    black := 1
    from := Point{0, 0}
    to := Point{1, 1}
    vulnerable := Vulnerable{Point{2, 2}, Point{3, 3}}
    enPassant := EnPassant{Point{4, 4}, Point{5, 5}}
    board := &MockBoard{}
    piece := Piece{white, 1}
    newPiece := Piece{white, 13}
    capturedPiece := Piece{black, 1}

	board.On("getEnPassant", white).Return(enPassant, nil)
    board.On("getVulnerable", white).Return(vulnerable, nil)
    promotionMove := createMoveSimple(board, piece, from, capturedPiece, to, newPiece)
    assert.Equal(t, promotionMove.allyDefense, false)

	board.On("setPiece", from, Piece{0, 0}).Return(true)
	board.On("setPiece", to, newPiece).Return(true)
	board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    err := promotionMove.execute()
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
    black := 0
    from := Point{0, 0}
    to := Point{1, 1}
    vulnerable := Vulnerable{Point{2, 2}, Point{3, 3}}
    enPassant := EnPassant{Point{4, 4}, Point{5, 5}}
    board := &MockBoard{}
    piece := Piece{white, 1}
    newPiece := Piece{white, 5}
    capturedPiece := Piece{black, 1}

    newEnPassant := EnPassant{Point{6, 6}, to}

	board.On("getEnPassant", white).Return(enPassant, nil)
    board.On("getVulnerable", white).Return(vulnerable, nil)
    revealEnPassantMove := createMoveRevealEnPassant(board, piece, from, capturedPiece, to, Piece{0, 0}, newEnPassant)
    assert.Equal(t, revealEnPassantMove.allyDefense, false)

	board.On("setPiece", from, Piece{0, 0}).Return(true)
	board.On("setPiece", to, newPiece).Return(true)
	board.On("setEnPassant", white, newEnPassant).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    err := revealEnPassantMove.execute()
    assert.Nil(t, err)

	board.On("setPiece", from, piece).Return(true)
	board.On("setPiece", to, capturedPiece).Return(true)
	board.On("setEnPassant", white, enPassant).Return(nil)
    board.On("setVulnerable", white, vulnerable).Return(nil)
    err = revealEnPassantMove.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
}

func Test_MoveCaptureEnPassant(t *testing.T) {
    white := 0
    black := 0
    from := Point{0, 0}
    to := Point{1, 1}
    vulnerable := Vulnerable{Point{2, 2}, Point{3, 3}}
    enPassant := EnPassant{Point{4, 4}, Point{5, 5}}
    board := &MockBoard{}
    piece := Piece{white, 1}
    newPiece := Piece{white, 5}
    capturedPiece := Piece{black, 1}

    enPassantRisk1 := Point{7, 7}
    enPassantRisk2 := Point{9, 9}
    enPassant1 := EnPassant{Point{6, 6}, enPassantRisk1}
    enPassant2 := EnPassant{Point{8, 8}, enPassantRisk2}
    enPassantCapturedPiece1 := Piece{black, 2}
    enPassantCapturedPiece2 := Piece{black, 3}

	board.On("getEnPassant", white).Return(enPassant, nil)
    board.On("getVulnerable", white).Return(vulnerable, nil)
    board.On("getPiece", enPassantRisk1).Return(enPassantCapturedPiece1, true)
    board.On("getPiece", enPassantRisk2).Return(enPassantCapturedPiece2, true)
    captureEnPassantMove := createMoveCaptureEnPassant(board, piece, from, capturedPiece, to, Piece{0, 0}, []EnPassant{enPassant1, enPassant2})
    assert.Equal(t, captureEnPassantMove.allyDefense, false)

	board.On("setPiece", from, Piece{0, 0}).Return(true)
	board.On("setPiece", to, newPiece).Return(true)
	board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
	board.On("setPiece", enPassantRisk1, Piece{0, 0}).Return(true)
	board.On("setPiece", enPassantRisk2, Piece{0, 0}).Return(true)
    err := captureEnPassantMove.execute()
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
}

func Test_MoveAllyDefense(t *testing.T) {
    white := 0
    from := Point{0, 0}
    to := Point{1, 1}
    piece := Piece{white, 1}
    board := &MockBoard{}

    allyDefenseMove := createMoveAllyDefense(board, piece, from, to)
    assert.Equal(t, allyDefenseMove.allyDefense, true)

    board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, Vulnerable{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    err := allyDefenseMove.execute()
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
    king := Piece{white, 1}
    newKing := Piece{white, 5}
    rook := Piece{white, 2}
    newRook := Piece{white, 6}

	board.On("getEnPassant", white).Return(enPassant, nil)
    board.On("getVulnerable", white).Return(vulnerable, nil)
    castleMove := createMoveCastle(board, king, from, rook, to, toKing, toRook, newVulnerable)
    assert.Equal(t, castleMove.allyDefense, false)

	board.On("setPiece", from, Piece{0, 0}).Return(true)
	board.On("setPiece", to, Piece{0, 0}).Return(true)
	board.On("setPiece", toKing, newKing).Return(true)
	board.On("setPiece", toRook, newRook).Return(true)
	board.On("setEnPassant", white, EnPassant{Point{-1, -1}, Point{-1, -1}}).Return(nil)
    board.On("setVulnerable", white, newVulnerable).Return(nil)
    err := castleMove.execute()
    assert.Nil(t, err)

	board.On("setPiece", toKing, Piece{0, 0}).Return(true)
	board.On("setPiece", toRook, Piece{0, 0}).Return(true)
	board.On("setPiece", from, king).Return(true)
	board.On("setPiece", to, rook).Return(true)
	board.On("setEnPassant", white, enPassant).Return(nil)
    board.On("setVulnerable", white, vulnerable).Return(nil)
    err = castleMove.undo()
    assert.Nil(t, err)

	board.AssertExpectations(t)
}

