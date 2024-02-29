package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMoveFactory struct {
	mock.Mock
}

func (m *MockMoveFactory) newSimpleMove(b Board, fromLocation Point, toLocation Point) (*SimpleMove, error) {
	args := m.Called(b, fromLocation, toLocation)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*SimpleMove), args.Error(1)
	}
}

func (m *MockMoveFactory) newRevealEnPassantMove(b Board, fromLocation Point, toLocation Point, target Point) (*RevealEnPassantMove, error) {
    args := m.Called(b, fromLocation, toLocation, target)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*RevealEnPassantMove), args.Error(1)
	}
}

func (m *MockMoveFactory) newCaptureEnPassantMove(b Board, fromLocation Point, toLocation Point) (*CaptureEnPassantMove, error) {
    args := m.Called(b, fromLocation, toLocation)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*CaptureEnPassantMove), args.Error(1)
	}
}

func (m *MockMoveFactory) newCastleMove(b Board, fromLocation Point, toLocation Point, toKingLocation Point, toRookLocation Point, newVulnerables []Point) (*CastleMove, error) {
    args := m.Called(b, fromLocation, toLocation, toKingLocation, toRookLocation, newVulnerables)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*CastleMove), args.Error(1)
	}
}

func (m *MockMoveFactory) newPromotionMove(move Move) (*PromotionMove, error) {
    args := m.Called(move)

    if args.Get(0) == nil {
        return nil, args.Error(1)
    } else {
        return args.Get(0).(*PromotionMove), args.Error(1)
    }
}

func (m *MockMoveFactory) newAllyDefenseMove(b Board, fromLocation Point, toLocation Point) (*AllyDefenseMove, error) {
    args := m.Called(b, fromLocation, toLocation)

    if args.Get(0) == nil {
        return nil, args.Error(1)
    } else {
        return args.Get(0).(*AllyDefenseMove), args.Error(1)
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

func (m *MockMove) getAction() Action {
    args := m.Called()
    return args.Get(0).(Action)
}

func (m *MockMove) getNewPiece() Piece {
    args := m.Called()
    return args.Get(0).(Piece)
}

func Test_SimpleMove(t *testing.T) {
	board := &MockBoard{}
	piece := &MockPiece{}
	newPiece := &MockPiece{}
	capturedPiece := &MockPiece{}
	en := EnPassant{}

	board.On("getPiece", Point{0, 0}).Return(piece, true)
	piece.On("copy").Return(newPiece)
    newPiece.On("setMoved").Return(nil)
	board.On("getPiece", Point{1, 1}).Return(capturedPiece, true)
	board.On("getEnPassant", "white").Return(en, nil)
	piece.On("getColor").Return("white")
    board.On("getVulnerables", "white").Return([]Point{{2, 2}})
	simpleMove, err := moveFactoryInstance.newSimpleMove(board, Point{0, 0}, Point{1, 1})
	assert.Nil(t, err)

	board.On("setPiece", Point{0, 0}, nil).Return(true)
	board.On("setPiece", Point{1, 1}, newPiece).Return(true)
	board.On("clearEnPassant", "white").Return(nil)
    board.On("setVulnerables", "white", []Point{}).Return(nil)
	err = simpleMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", Point{0, 0}, piece).Return(true)
	board.On("setPiece", Point{1, 1}, capturedPiece).Return(true)
	board.On("setEnPassant", "white", en).Return(nil)
    board.On("setVulnerables", "white", []Point{{2, 2}}).Return(nil)
	err = simpleMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
    newPiece.AssertExpectations(t)
}

func Test_RevealEnPassantMove(t *testing.T) {
	board := &MockBoard{}
	piece := &MockPiece{}
	newPiece := &MockPiece{}
	capturedPiece := &MockPiece{}
	en := EnPassant{}

	board.On("getPiece", Point{0, 0}).Return(piece, true)
	piece.On("copy").Return(newPiece)
    newPiece.On("setMoved").Return(nil)
	board.On("getPiece", Point{2, 2}).Return(capturedPiece, true)
	board.On("getEnPassant", "white").Return(en, nil)
	piece.On("getColor").Return("white")
    board.On("getVulnerables", "white").Return([]Point{{2, 2}})
	revealEnPassantMove, err := moveFactoryInstance.newRevealEnPassantMove(board, Point{0, 0}, Point{2, 2}, Point{1, 1})
	assert.Nil(t, err)

	board.On("setPiece", Point{0, 0}, nil).Return(true)
	board.On("setPiece", Point{2, 2}, newPiece).Return(true)
	board.On("setEnPassant", "white", EnPassant{Point{1, 1}, Point{2, 2}}).Return(nil)
    board.On("setVulnerables", "white", []Point{}).Return(nil)
	err = revealEnPassantMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", Point{0, 0}, piece).Return(true)
	board.On("setPiece", Point{2, 2}, capturedPiece).Return(true)
	board.On("setEnPassant", "white", en).Return(nil)
    board.On("setVulnerables", "white", []Point{{2, 2}}).Return(nil)
	err = revealEnPassantMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
    newPiece.AssertExpectations(t)
}

func Test_CaptureEnPassantMove(t *testing.T) {
	board := &MockBoard{}
	piece := &MockPiece{}
	newPiece := &MockPiece{}
	capturedPiece := &MockPiece{}
	en := EnPassant{}
	encPiece := &MockPiece{}

	board.On("getPiece", Point{0, 0}).Return(piece, true)
    piece.On("copy").Return(newPiece)
    newPiece.On("setMoved").Return(nil)
	board.On("getPiece", Point{1, 1}).Return(capturedPiece, true)
	board.On("getEnPassant", "white").Return(en, nil)
	piece.On("getColor").Return("white")
	board.On("possibleEnPassant", "white", Point{1, 1}).Return([]EnPassant{{Point{1, 1}, Point{2, 2}}}, nil)
	board.On("getPiece", Point{2, 2}).Return(encPiece, true)
    board.On("getVulnerables", "white").Return([]Point{{2, 2}})
	captureEnPassantMove, err := moveFactoryInstance.newCaptureEnPassantMove(board, Point{0, 0}, Point{1, 1})
	assert.Nil(t, err)

	assert.NotNil(t, captureEnPassantMove)

	board.On("setPiece", Point{0, 0}, nil).Return(true)
	board.On("setPiece", Point{1, 1}, newPiece).Return(true)
	board.On("setPiece", Point{2, 2}, nil).Return(true)
	board.On("clearEnPassant", "white").Return(nil)
    board.On("setVulnerables", "white", []Point{}).Return(nil)
	err = captureEnPassantMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", Point{0, 0}, piece).Return(true)
	board.On("setPiece", Point{1, 1}, capturedPiece).Return(true)
	board.On("setPiece", Point{2, 2}, encPiece).Return(true)
	board.On("setEnPassant", "white", en).Return(nil)
    board.On("setVulnerables", "white", []Point{{2, 2}}).Return(nil)
	err = captureEnPassantMove.undo()
	assert.Nil(t, err)

    board.AssertExpectations(t)
    piece.AssertExpectations(t)
    newPiece.AssertExpectations(t)
}

func Test_CastleMove(t *testing.T) {
	board := &MockBoard{}
	king := &MockPiece{}
	newKing := &MockPiece{}
	rook := &MockPiece{}
	newRook := &MockPiece{}
	en := EnPassant{}
    newVulnerables := []Point{{4, 4}}

	board.On("getPiece", Point{0, 0}).Return(king, true)
	board.On("getPiece", Point{1, 1}).Return(rook, true)
	king.On("copy").Return(newKing)
	rook.On("copy").Return(newRook)
    newKing.On("setMoved").Return(nil)
    newRook.On("setMoved").Return(nil)
	board.On("getEnPassant", "white").Return(en, nil)
	king.On("getColor").Return("white")
    board.On("getVulnerables", "white").Return([]Point{{5, 5}})
	castleMove, err := moveFactoryInstance.newCastleMove(board, Point{0, 0}, Point{1, 1}, Point{2, 2}, Point{3, 3}, newVulnerables)
	assert.Nil(t, err)

	board.On("setPiece", Point{0, 0}, nil).Return(true)
	board.On("setPiece", Point{1, 1}, nil).Return(true)
	board.On("setPiece", Point{2, 2}, newKing).Return(true)
	board.On("setPiece", Point{3, 3}, newRook).Return(true)
	board.On("clearEnPassant", "white").Return(nil)
    board.On("setVulnerables", "white", []Point{{4, 4}}).Return(nil)
	err = castleMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", Point{0, 0}, king).Return(true)
	board.On("setPiece", Point{1, 1}, rook).Return(true)
	board.On("setPiece", Point{2, 2}, nil).Return(true)
	board.On("setPiece", Point{3, 3}, nil).Return(true)
	board.On("setEnPassant", "white", en).Return(nil)
    board.On("setVulnerables", "white", []Point{{5, 5}}).Return(nil)
	err = castleMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	king.AssertExpectations(t)
	rook.AssertExpectations(t)
    newKing.AssertExpectations(t)
    newRook.AssertExpectations(t)
}

func Test_PromotionMove(t *testing.T) {
	board := &MockBoard{}
	piece := &MockPiece{}
	newPiece := &MockPiece{}
	capturedPiece := &MockPiece{}
	en := EnPassant{}
    queen := &MockPiece{}

	board.On("getPiece", Point{0, 0}).Return(piece, true)
	piece.On("copy").Return(newPiece)
    newPiece.On("setMoved").Return(nil)
	board.On("getPiece", Point{1, 1}).Return(capturedPiece, true)
	board.On("getEnPassant", "white").Return(en, nil)
	piece.On("getColor").Return("white")
    board.On("getVulnerables", "white").Return([]Point{{2, 2}})
	simpleMove, err := moveFactoryInstance.newSimpleMove(board, Point{0, 0}, Point{1, 1})
	assert.Nil(t, err)
    promotionMove, err := moveFactoryInstance.newPromotionMove(simpleMove)
    assert.Nil(t, err)
    err = promotionMove.setPromotionPiece(queen)
    assert.Nil(t, err)

	board.On("setPiece", Point{0, 0}, nil).Return(true)
	board.On("setPiece", Point{1, 1}, newPiece).Return(true)
    board.On("setPiece", Point{1, 1}, queen).Return(true)
	board.On("clearEnPassant", "white").Return(nil)
    board.On("setVulnerables", "white", []Point{}).Return(nil)
	err = promotionMove.execute()
	assert.Nil(t, err)

	board.On("setPiece", Point{0, 0}, piece).Return(true)
    board.On("setPiece", Point{1, 1}, newPiece).Return(true)
	board.On("setPiece", Point{1, 1}, capturedPiece).Return(true)
	board.On("setEnPassant", "white", en).Return(nil)
    board.On("setVulnerables", "white", []Point{{2, 2}}).Return(nil)
	err = promotionMove.undo()
	assert.Nil(t, err)

	board.AssertExpectations(t)
	piece.AssertExpectations(t)
    newPiece.AssertExpectations(t)
}

