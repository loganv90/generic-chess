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

    testPiece1 := Piece{}
    testPiece2 := Piece{}
    testEnPassant := EnPassant{}
    testVulnerable := Vulnerable{}

    moves := Array100[FastMove]{}
    move := moves.get()

	board.On("getEnPassant", white).Return(&enPassant).Once()
    board.On("getVulnerable", white).Return(&vulnerable).Once()
    addMoveSimple(board, &piece, &from, &capturedPiece, &to, nil, &moves)
    assert.Equal(t, move.allyDefense, false)

	board.On("getPiece", &from).Return(&testPiece1).Once()
	board.On("getPiece", &to).Return(&testPiece2).Once()
	board.On("getEnPassant", white).Return(&testEnPassant).Once()
    board.On("getVulnerable", white).Return(&testVulnerable).Once()
    err := move.execute()
    assert.Nil(t, err)
    assert.Equal(t, testPiece1, Piece{0, 0})
    assert.Equal(t, testPiece2, newPiece)
    assert.Equal(t, testEnPassant, EnPassant{Point{-1, -1}, Point{-1, -1}})
    assert.Equal(t, testVulnerable, Vulnerable{Point{-1, -1}, Point{-1, -1}})

	board.On("getPiece", &from).Return(&testPiece1).Once()
	board.On("getPiece", &to).Return(&testPiece2).Once()
	board.On("getEnPassant", white).Return(&testEnPassant).Once()
    board.On("getVulnerable", white).Return(&testVulnerable).Once()
    err = move.undo()
    assert.Nil(t, err)
    assert.Equal(t, testPiece1, piece)
    assert.Equal(t, testPiece2, capturedPiece)
    assert.Equal(t, testEnPassant, enPassant)
    assert.Equal(t, testVulnerable, vulnerable)

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

    testPiece1 := Piece{}
    testPiece2 := Piece{}
    testEnPassant := EnPassant{}
    testVulnerable := Vulnerable{}

    moves := Array100[FastMove]{}
    move := moves.get()

	board.On("getEnPassant", white).Return(&enPassant).Once()
    board.On("getVulnerable", white).Return(&vulnerable).Once()
    addMoveSimple(board, &piece, &from, &capturedPiece, &to, &newPiece, &moves)
    assert.Equal(t, move.allyDefense, false)

	board.On("getPiece", &from).Return(&testPiece1).Once()
	board.On("getPiece", &to).Return(&testPiece2).Once()
	board.On("getEnPassant", white).Return(&testEnPassant).Once()
    board.On("getVulnerable", white).Return(&testVulnerable).Once()
    err := move.execute()
    assert.Nil(t, err)
    assert.Equal(t, testPiece1, Piece{0, 0})
    assert.Equal(t, testPiece2, newPiece)
    assert.Equal(t, testEnPassant, EnPassant{Point{-1, -1}, Point{-1, -1}})
    assert.Equal(t, testVulnerable, Vulnerable{Point{-1, -1}, Point{-1, -1}})

	board.On("getPiece", &from).Return(&testPiece1).Once()
	board.On("getPiece", &to).Return(&testPiece2).Once()
	board.On("getEnPassant", white).Return(&testEnPassant).Once()
    board.On("getVulnerable", white).Return(&testVulnerable).Once()
    err = move.undo()
    assert.Nil(t, err)
    assert.Equal(t, testPiece1, piece)
    assert.Equal(t, testPiece2, capturedPiece)
    assert.Equal(t, testEnPassant, enPassant)
    assert.Equal(t, testVulnerable, vulnerable)

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

    testPiece1 := Piece{}
    testPiece2 := Piece{}
    testEnPassant := EnPassant{}
    testVulnerable := Vulnerable{}

    moves := Array100[FastMove]{}
    move := moves.get()

	board.On("getEnPassant", white).Return(&enPassant).Once()
    board.On("getVulnerable", white).Return(&vulnerable).Once()
    addMoveRevealEnPassant(board, &piece, &from, &capturedPiece, &to, nil, &newEnPassant, &moves)
    assert.Equal(t, move.allyDefense, false)

	board.On("getPiece", &from).Return(&testPiece1).Once()
	board.On("getPiece", &to).Return(&testPiece2).Once()
	board.On("getEnPassant", white).Return(&testEnPassant).Once()
    board.On("getVulnerable", white).Return(&testVulnerable).Once()
    err := move.execute()
    assert.Nil(t, err)
    assert.Equal(t, testPiece1, Piece{0, 0})
    assert.Equal(t, testPiece2, newPiece)
    assert.Equal(t, testEnPassant, newEnPassant)
    assert.Equal(t, testVulnerable, Vulnerable{Point{-1, -1}, Point{-1, -1}})

	board.On("getPiece", &from).Return(&testPiece1).Once()
	board.On("getPiece", &to).Return(&testPiece2).Once()
	board.On("getEnPassant", white).Return(&testEnPassant).Once()
    board.On("getVulnerable", white).Return(&testVulnerable).Once()
    err = move.undo()
    assert.Nil(t, err)
    assert.Equal(t, testPiece1, piece)
    assert.Equal(t, testPiece2, capturedPiece)
    assert.Equal(t, testEnPassant, enPassant)
    assert.Equal(t, testVulnerable, vulnerable)

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

    testPiece1 := Piece{}
    testPiece2 := Piece{}
    testPiece3 := Piece{}
    testPiece4 := Piece{}
    testEnPassant := EnPassant{}
    testVulnerable := Vulnerable{}

    moves := Array100[FastMove]{}
    move := moves.get()

	board.On("getEnPassant", white).Return(&enPassant).Once()
    board.On("getVulnerable", white).Return(&vulnerable).Once()
    board.On("getPiece", &enPassantRisk1).Return(&enPassantCapturedPiece1).Once()
    board.On("getPiece", &enPassantRisk2).Return(&enPassantCapturedPiece2).Once()
    addMoveCaptureEnPassant(board, &piece, &from, &capturedPiece, &to, nil, &moves, &enPassant1, &enPassant2)
    assert.Equal(t, move.allyDefense, false)

	board.On("getPiece", &from).Return(&testPiece1).Once()
	board.On("getPiece", &to).Return(&testPiece2).Once()
	board.On("getEnPassant", white).Return(&testEnPassant).Once()
    board.On("getVulnerable", white).Return(&testVulnerable).Once()
	board.On("getPiece", &enPassantRisk1).Return(&testPiece3).Once()
	board.On("getPiece", &enPassantRisk2).Return(&testPiece4).Once()
    err := move.execute()
    assert.Nil(t, err)
    assert.Equal(t, testPiece1, Piece{0, 0})
    assert.Equal(t, testPiece2, newPiece)
    assert.Equal(t, testEnPassant, EnPassant{Point{-1, -1}, Point{-1, -1}})
    assert.Equal(t, testVulnerable, Vulnerable{Point{-1, -1}, Point{-1, -1}})
    assert.Equal(t, testPiece3, Piece{0, 0})
    assert.Equal(t, testPiece4, Piece{0, 0})

	board.On("getPiece", &from).Return(&testPiece1).Once()
	board.On("getPiece", &to).Return(&testPiece2).Once()
	board.On("getEnPassant", white).Return(&testEnPassant).Once()
    board.On("getVulnerable", white).Return(&testVulnerable).Once()
	board.On("getPiece", &enPassantRisk1).Return(&testPiece3).Once()
	board.On("getPiece", &enPassantRisk2).Return(&testPiece4).Once()
    err = move.undo()
    assert.Nil(t, err)
    assert.Equal(t, testPiece1, piece)
    assert.Equal(t, testPiece2, capturedPiece)
    assert.Equal(t, testEnPassant, enPassant)
    assert.Equal(t, testVulnerable, vulnerable)
    assert.Equal(t, testPiece3, enPassantCapturedPiece1)
    assert.Equal(t, testPiece4, enPassantCapturedPiece2)

	board.AssertExpectations(t)
}

func Test_MoveAllyDefense(t *testing.T) {
    white := 0
    from := Point{0, 0}
    to := Point{1, 1}
    piece := Piece{white, 1}
    board := &MockBoard{}

    moves := Array100[FastMove]{}
    move := moves.get()

    addMoveAllyDefense(board, &piece, &from, &to, &moves)
    assert.Equal(t, move.allyDefense, true)

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

    testPiece1 := Piece{}
    testPiece2 := Piece{}
    testPiece3 := Piece{}
    testPiece4 := Piece{}
    testEnPassant := EnPassant{}
    testVulnerable := Vulnerable{}

    moves := Array100[FastMove]{}
    move := moves.get()

	board.On("getEnPassant", white).Return(&enPassant).Once()
    board.On("getVulnerable", white).Return(&vulnerable).Once()
    addMoveCastle(board, &king, &from, &toKing, &rook, &to, &toRook, &newVulnerable, &moves)
    assert.Equal(t, move.allyDefense, false)

	board.On("getPiece", &from).Return(&testPiece1).Once()
	board.On("getPiece", &to).Return(&testPiece2).Once()
	board.On("getPiece", &toKing).Return(&testPiece3).Once()
	board.On("getPiece", &toRook).Return(&testPiece4).Once()
	board.On("getEnPassant", white).Return(&testEnPassant).Once()
    board.On("getVulnerable", white).Return(&testVulnerable).Once()
    err := move.execute()
    assert.Nil(t, err)
    assert.Equal(t, testPiece1, Piece{0, 0})
    assert.Equal(t, testPiece2, Piece{0, 0})
    assert.Equal(t, testPiece3, newKing)
    assert.Equal(t, testPiece4, newRook)
    assert.Equal(t, testEnPassant, EnPassant{Point{-1, -1}, Point{-1, -1}})
    assert.Equal(t, testVulnerable, newVulnerable)

	board.On("getPiece", &from).Return(&testPiece1).Once()
	board.On("getPiece", &to).Return(&testPiece2).Once()
	board.On("getPiece", &toKing).Return(&testPiece3).Once()
	board.On("getPiece", &toRook).Return(&testPiece4).Once()
	board.On("getEnPassant", white).Return(&testEnPassant).Once()
    board.On("getVulnerable", white).Return(&testVulnerable).Once()
    err = move.undo()
    assert.Nil(t, err)
    assert.Equal(t, testPiece1, king)
    assert.Equal(t, testPiece2, rook)
    assert.Equal(t, testPiece3, Piece{0, 0})
    assert.Equal(t, testPiece4, Piece{0, 0})
    assert.Equal(t, testEnPassant, enPassant)
    assert.Equal(t, testVulnerable, vulnerable)

	board.AssertExpectations(t)
}

