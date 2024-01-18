package chess


import (
	"testing"
    "slices"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBoard struct {
	mock.Mock
}

func (m *MockBoard) getPiece(location *Point) (Piece, error) {
	args := m.Called(location)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(Piece), args.Error(1)
	}
}

func (m *MockBoard) setPiece(location *Point, piece Piece) error {
	args := m.Called(location, piece)
	return args.Error(0)
}

func (m *MockBoard) disableLocation(location *Point) error {
    args := m.Called(location)
    return args.Error(0)
}

func (m *MockBoard) getVulnerables(color string) ([]*Point, error) {
    args := m.Called(color)
    return args.Get(0).([]*Point), args.Error(1)
}

func (m *MockBoard) setVulnerables(color string, vulnerables []*Point) error {
    args := m.Called(color, vulnerables)
    return args.Error(0)
}

func (m *MockBoard) getEnPassant(color string) (*EnPassant, error) {
	args := m.Called(color)
	return args.Get(0).(*EnPassant), args.Error(1)
}

func (m *MockBoard) setEnPassant(color string, enPassant *EnPassant) error {
    args := m.Called(color, enPassant)
    return args.Error(0)
}

func (m *MockBoard) possibleEnPassant(color string, location *Point) ([]*EnPassant, error) {
    args := m.Called(color, location)
    return args.Get(0).([]*EnPassant), args.Error(1)
}

func (m *MockBoard) clearEnPassant(color string) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockBoard) PotentialMoves(fromLocation *Point) ([]Move, error) {
    args := m.Called(fromLocation)
    return args.Get(0).([]Move), args.Error(1)
}

func (m *MockBoard) ValidMoves(fromLocation *Point) ([]Move, error) {
    args := m.Called(fromLocation)
    return args.Get(0).([]Move), args.Error(1)
}

func (m *MockBoard) CalculateMoves(color string) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockBoard) Size() *Point {
    args := m.Called()
    return args.Get(0).(*Point)
}

func (m *MockBoard) Print() string {
    args := m.Called()
    return args.String(0)
}

func (m *MockBoard) State() *BoardData {
    args := m.Called()
    return args.Get(0).(*BoardData)
}

func (m *MockBoard) Checkmate() bool {
    args := m.Called()
    return args.Bool(0)
}

func (m *MockBoard) Stalemate() bool {
    args := m.Called()
    return args.Bool(0)
}

func Test_NewSimpleBoard_DefaultFen(t *testing.T) {
    s, err := createSimpleBoardWithDefaultPieceLocations()
	assert.Nil(t, err)

	for y := 2; y <= 5; y++ {
		for x := 0; x <= 7; x++ {
			piece, err := s.getPiece(&Point{x, y})
			assert.Nil(t, err)
			assert.Nil(t, piece)
		}
	}

	for _, y := range []int{1, 6} {
		for x := 0; x <= 7; x++ {
			piece, err := s.getPiece(&Point{x, y})
			assert.Nil(t, err)
			_, ok := piece.(*Pawn)
			assert.True(t, ok)
		}
	}

	for _, y := range []int{0, 7} {
		piece, err := s.getPiece(&Point{0, y})
		assert.Nil(t, err)
		_, ok := piece.(*Rook)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{1, y})
		assert.Nil(t, err)
		_, ok = piece.(*Knight)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{2, y})
		assert.Nil(t, err)
		_, ok = piece.(*Bishop)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{3, y})
		assert.Nil(t, err)
		_, ok = piece.(*Queen)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{4, y})
		assert.Nil(t, err)
		_, ok = piece.(*King)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{5, y})
		assert.Nil(t, err)
		_, ok = piece.(*Bishop)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{6, y})
		assert.Nil(t, err)
		_, ok = piece.(*Knight)
		assert.True(t, ok)

		piece, err = s.getPiece(&Point{7, y})
		assert.Nil(t, err)
		_, ok = piece.(*Rook)
		assert.True(t, ok)
	}
}

func Test_getAndSetPiece(t *testing.T) {
    s, err := newSimpleBoard(&Point{8, 8})
    assert.Nil(t, err)

    p := newPawn("white", false, 0, 1)
    location := &Point{0, 0}

    // adding piece to board
    err = s.setPiece(location, p)
    assert.Nil(t, err)

    piece, err := s.getPiece(location)
    assert.Nil(t, err)
    assert.Equal(t, p, piece)
    pieceLocations, ok := s.pieceLocationsMap["white"]
    assert.True(t, ok)
    exists := slices.Contains(pieceLocations, location)
    assert.True(t, exists)

    // removing piece from board
    err = s.setPiece(location, nil)
    assert.Nil(t, err)

    piece, err = s.getPiece(location)
    assert.Nil(t, err)
    assert.Nil(t, piece)
    pieceLocations, ok = s.pieceLocationsMap["white"]
    assert.True(t, ok)
    exists = slices.Contains(pieceLocations, location)
    assert.False(t, exists)
}

func Test_CalculateMoves_default(t *testing.T) {
    s, err := createSimpleBoardWithDefaultPieceLocations()
    assert.Nil(t, err)

    s.CalculateMoves("white")
    boardData := s.State()

    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 20, moveCount)
    assert.False(t, boardData.Check)
    assert.False(t, boardData.Checkmate)
    assert.False(t, boardData.Stalemate)
}

func Test_CalculateMoves_check(t *testing.T) {
    s, err := newSimpleBoard(&Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(&Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(&Point{0, 1}, newQueen("black"))
    s.setPiece(&Point{0, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves("white")
    boardData := s.State()

    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 1, moveCount)
    assert.True(t, boardData.Check)
    assert.False(t, boardData.Checkmate)
    assert.False(t, boardData.Stalemate)
}

func Test_CalculateMoves_checkmate(t *testing.T) {
    s, err := newSimpleBoard(&Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(&Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(&Point{0, 1}, newQueen("black"))
    s.setPiece(&Point{0, 2}, newKing("black", false, 0, -1))

    s.CalculateMoves("white")
    boardData := s.State()

    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 0, moveCount)
    assert.True(t, boardData.Check)
    assert.True(t, boardData.Checkmate)
    assert.False(t, boardData.Stalemate)
}

func Test_CalculateMoves_stalemate(t *testing.T) {
    s, err := newSimpleBoard(&Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(&Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(&Point{1, 2}, newQueen("black"))
    s.setPiece(&Point{0, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves("white")
    boardData := s.State()

    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 0, moveCount)
    assert.False(t, boardData.Check)
    assert.False(t, boardData.Checkmate)
    assert.True(t, boardData.Stalemate)
}

func Test_CalculateMoves_noCastleThroughCheck(t *testing.T) {
    s, err := newSimpleBoard(&Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(&Point{4, 0}, newKing("white", false, 0, 1))
    s.setPiece(&Point{0, 0}, newRook("white", false))
    s.setPiece(&Point{3, 7}, newRook("black", false))
    s.setPiece(&Point{4, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves("white")
    boardData := s.State()

    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 13, moveCount)
    assert.False(t, boardData.Check)
    assert.False(t, boardData.Checkmate)
    assert.False(t, boardData.Stalemate)
}

func Test_CalculateMoves_promotion(t *testing.T) {
    s, err := newSimpleBoard(&Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(&Point{7, 6}, newPawn("white", true, 0, 1))
    s.setPiece(&Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(&Point{0, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves("white")
    boardData := s.State()

    moveCount := 0
    for _, m := range s.moveMap {
        moveCount += len(m)
    }

    assert.Equal(t, 7, moveCount)
    assert.False(t, boardData.Check)
    assert.False(t, boardData.Checkmate)
    assert.False(t, boardData.Stalemate)
}

