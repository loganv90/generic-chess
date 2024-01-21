package chess

import (
	"testing"
    "errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPiece struct {
	mock.Mock
}

func (m *MockPiece) getColor() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockPiece) copy() Piece {
    args := m.Called()
    return args.Get(0).(Piece)
}

func (m *MockPiece) setMoved() error {
    args := m.Called()
    return args.Error(0)
}

func (m *MockPiece) moves(board Board, location *Point) []Move {
    args := m.Called(board, location)
	return args.Get(0).([]Move)
}

func (m *MockPiece) print() string {
	args := m.Called()
	return args.String(0)
}

func Test_Pawn_Moves_Unmoved(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
	board.On("getPiece", mock.Anything, mock.Anything).Return(nil, nil)
	board.On("possibleEnPassant", mock.Anything, mock.Anything).Return([]*EnPassant{}, nil)

	moveFactory := &MockMoveFactory{}
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{3, 4}).Return(nil, nil)
	moveFactory.On("newRevealEnPassantMove", board, &Point{3, 3}, &Point{3, 5}, &Point{3, 4}).Return(nil, nil)
	moveFactoryInstance = moveFactory

	pawn := newPawn("white", false, 0, 1)

	moves := pawn.moves(board, &Point{3, 3})
	assert.Len(t, moves, 2)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func Test_Pawn_Moves_Moved(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
	board.On("getPiece", mock.Anything, mock.Anything).Return(nil, nil)
	board.On("possibleEnPassant", mock.Anything, mock.Anything).Return([]*EnPassant{}, nil)

	moveFactory := &MockMoveFactory{}
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{3, 4}).Return(nil, nil)
	moveFactoryInstance = moveFactory

	pawn := newPawn("white", true, 0, 1)

	moves := pawn.moves(board, &Point{3, 3})
	assert.Len(t, moves, 1)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func Test_Pawn_Moves_Capturing(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
	board.On("getPiece", &Point{2, 4}).Return(nil, nil)
	board.On("getPiece", &Point{3, 4}).Return(nil, nil)
	board.On("getPiece", &Point{3, 5}).Return(nil, nil)
    board.On("getPiece", &Point{3, 6}).Return(nil, nil)
    board.On("getPiece", &Point{4, 5}).Return(nil, nil)
	board.On("possibleEnPassant", mock.Anything, mock.Anything).Return([]*EnPassant{}, nil)

	moveFactory := &MockMoveFactory{}
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{3, 4}).Return(nil, nil)
	moveFactory.On("newRevealEnPassantMove", board, &Point{3, 3}, &Point{3, 5}, &Point{3, 4}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{4, 4}).Return(nil, nil)
	moveFactoryInstance = moveFactory

	pawn := newPawn("white", false, 0, 1)

	blackPawn := newPawn("black", false, 0, 1)
	board.On("getPiece", &Point{4, 4}).Return(blackPawn, nil)

	moves := pawn.moves(board, &Point{3, 3})
	assert.Len(t, moves, 3)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func Test_Pawn_Moves_CapturingEnPassant(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
	board.On("getPiece", &Point{2, 4}).Return(nil, nil)
	board.On("getPiece", &Point{3, 4}).Return(nil, nil)
	board.On("getPiece", &Point{3, 5}).Return(nil, nil)
	board.On("getPiece", &Point{4, 4}).Return(nil, nil)
    board.On("getPiece", &Point{3, 6}).Return(nil, nil)
    board.On("getPiece", &Point{4, 5}).Return(nil, nil)
	board.On("possibleEnPassant", "white", &Point{4, 4}).Return([]*EnPassant{{&Point{4, 4}, &Point{3, 4}}}, nil)
	board.On("possibleEnPassant", "white", &Point{2, 4}).Return([]*EnPassant{}, nil)

	moveFactory := &MockMoveFactory{}
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{3, 4}).Return(nil, nil)
	moveFactory.On("newRevealEnPassantMove", board, &Point{3, 3}, &Point{3, 5}, &Point{3, 4}).Return(nil, nil)
	moveFactory.On("newCaptureEnPassantMove", board, &Point{3, 3}, &Point{4, 4}).Return(nil, nil)
	moveFactoryInstance = moveFactory

	pawn := newPawn("white", false, 0, 1)

	moves := pawn.moves(board, &Point{3, 3})
	assert.Len(t, moves, 3)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func Test_Pawn_Moves_Promotion(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
	board.On("getPiece", &Point{3, 4}).Return(nil, nil)
	board.On("getPiece", &Point{3, 5}).Return(nil, errors.New(""))
    board.On("getPiece", &Point{4, 4}).Return(nil, nil)
    board.On("getPiece", &Point{2, 4}).Return(nil, nil)
	board.On("possibleEnPassant", mock.Anything, mock.Anything).Return([]*EnPassant{}, nil)

	moveFactory := &MockMoveFactory{}
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{3, 4}).Return(nil, nil)
    moveFactory.On("newPromotionMoves", mock.Anything, mock.Anything).Return([]*PromotionMove{nil, nil, nil}, nil)
	moveFactoryInstance = moveFactory

	pawn := newPawn("white", false, 0, 1)

	moves := pawn.moves(board, &Point{3, 3})
	assert.Len(t, moves, 3)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func Test_Knight_Moves(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
    moveFactory := &MockMoveFactory{}

    inPoints := []*Point{
        {2, 3},
        {3, 2},
        {0, 3},
        {3, 0},
    }

    outPoints := []*Point{
        {-1, 2},
        {2, -1},
        {0, -1},
        {-1, 0},
    }

    for _, p := range inPoints {
        board.On("getPiece", p).Return(nil, nil)
        moveFactory.On("newSimpleMove", board, &Point{1, 1}, p).Return(nil, nil)
    }

    for _, p := range outPoints {
        board.On("getPiece", p).Return(nil, errors.New(""))
    }

	moveFactoryInstance = moveFactory

	knight := newKnight("white")

	moves := knight.moves(board, &Point{1, 1})
	assert.Len(t, moves, 4)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func Test_Bishop_Moves(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
    moveFactory := &MockMoveFactory{}

    inPoints := []*Point{
        {2, 2},
        {3, 3},
        {4, 4},
        {5, 5},
        {6, 6},
        {7, 7},
        {0, 2},
        {2, 0},
        {0, 0},
    }

    outPoints := []*Point{
        {8, 8},
        {-1, 3},
        {3, -1},
        {-1, -1},
    }

    for _, p := range inPoints {
        board.On("getPiece", p).Return(nil, nil)
        moveFactory.On("newSimpleMove", board, &Point{1, 1}, p).Return(nil, nil)
    }

    for _, p := range outPoints {
        board.On("getPiece", p).Return(nil, errors.New(""))
    }

	moveFactoryInstance = moveFactory

	bishop := newBishop("white")

	moves := bishop.moves(board, &Point{1, 1})
	assert.Len(t, moves, 9)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func Test_Rook_Moves(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
	moveFactory := &MockMoveFactory{}

    inPoints := []*Point{
        {2, 1},
        {3, 1},
        {4, 1},
        {5, 1},
        {6, 1},
        {7, 1},
        {1, 2},
        {1, 3},
        {1, 4},
        {1, 5},
        {1, 6},
        {1, 7},
        {0, 1},
        {1, 0},
    }

    outPoints := []*Point{
        {8, 1},
        {1, 8},
        {1, -1},
        {-1, 1},
    }

    for _, p := range inPoints {
        board.On("getPiece", p).Return(nil, nil)
        moveFactory.On("newSimpleMove", board, &Point{1, 1}, p).Return(nil, nil)
    }

    for _, p := range outPoints {
        board.On("getPiece", p).Return(nil, errors.New(""))
    }

	moveFactoryInstance = moveFactory

	rook := newRook("white", false)

	moves := rook.moves(board, &Point{1, 1})
	assert.Len(t, moves, 14)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func Test_Queen_Moves(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
	moveFactory := &MockMoveFactory{}

    inPoints := []*Point{
        {2, 2},
        {3, 3},
        {4, 4},
        {5, 5},
        {6, 6},
        {7, 7},
        {0, 2},
        {2, 0},
        {0, 0},
        {2, 1},
        {3, 1},
        {4, 1},
        {5, 1},
        {6, 1},
        {7, 1},
        {1, 2},
        {1, 3},
        {1, 4},
        {1, 5},
        {1, 6},
        {1, 7},
        {0, 1},
        {1, 0},
    }

    outPoints := []*Point{
        {8, 8},
        {-1, 3},
        {3, -1},
        {-1, -1},
        {8, 1},
        {1, 8},
        {1, -1},
        {-1, 1},
    }

    for _, p := range inPoints {
        board.On("getPiece", p).Return(nil, nil)
        moveFactory.On("newSimpleMove", board, &Point{1, 1}, p).Return(nil, nil)
    }

    for _, p := range outPoints {
        board.On("getPiece", p).Return(nil, errors.New(""))
    }

	moveFactoryInstance = moveFactory

	queen := newQueen("white")

	moves := queen.moves(board, &Point{1, 1})
	assert.Len(t, moves, 23)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func Test_King_Moves_CanCastleAndUnmoved(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
	moveFactory := &MockMoveFactory{}

    inPoints := []*Point{
        {2, 2},
        {3, 2},
        {4, 2},
        {2, 3},
        {4, 3},
        {2, 4},
        {3, 4},
        {4, 4},
        {5, 3},
        {1, 3},
        {2, 3},
        {4, 3},
        {5, 3},
        {6, 3},
    }

    for _, p := range inPoints {
        board.On("getPiece", p).Return(nil, nil)
    }

	board.On("Size").Return(&Point{8, 8})

	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{2, 2}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{3, 2}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{4, 2}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{2, 3}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{4, 3}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{2, 4}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{3, 4}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{4, 4}).Return(nil, nil)
	moveFactory.On("newCastleMove", board, &Point{3, 3}, &Point{0, 3}, &Point{2, 3}, &Point{3, 3}, []*Point{}).Return(nil, nil)
	moveFactory.On("newCastleMove", board, &Point{3, 3}, &Point{7, 3}, &Point{6, 3}, &Point{5, 3}, []*Point{{4, 3}, {5, 3}}).Return(nil, nil)
	moveFactoryInstance = moveFactory

	king := newKing("white", false, 0, 1)

	rook := newRook("white", false)
	board.On("getPiece", &Point{0, 3}).Return(rook, nil)
	board.On("getPiece", &Point{7, 3}).Return(rook, nil)

	moves := king.moves(board, &Point{3, 3})
	assert.Len(t, moves, 10)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

func Test_King_Moves_CanCastleAndMoved(t *testing.T) {
    t.Cleanup(func() { moveFactoryInstance = &ConcreteMoveFactory{} })

	board := &MockBoard{}
	moveFactory := &MockMoveFactory{}

    inPoints := []*Point{
        {2, 2},
        {3, 2},
        {4, 2},
        {2, 3},
        {4, 3},
        {2, 4},
        {3, 4},
        {4, 4},
    }

    for _, p := range inPoints {
        board.On("getPiece", p).Return(nil, nil)
    }

	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{2, 2}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{3, 2}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{4, 2}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{2, 3}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{4, 3}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{2, 4}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{3, 4}).Return(nil, nil)
	moveFactory.On("newSimpleMove", board, &Point{3, 3}, &Point{4, 4}).Return(nil, nil)
	moveFactoryInstance = moveFactory

	king := newKing("white", true, 0, 1)

	moves := king.moves(board, &Point{3, 3})
	assert.Len(t, moves, 8)

	board.AssertExpectations(t)
	moveFactory.AssertExpectations(t)
}

