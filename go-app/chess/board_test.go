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

func (m *MockBoard) disablePieces(color string, disable bool) error {
    args := m.Called(color, disable)
    return args.Error(0)
}

func (m *MockBoard) getPiece(location Point) (Piece, bool) {
	args := m.Called(location)

	if args.Get(0) == nil {
		return nil, args.Bool(1)
	} else {
		return args.Get(0).(Piece), args.Bool(1)
	}
}

func (m *MockBoard) setPiece(location Point, piece Piece) bool {
	args := m.Called(location, piece)
    return args.Bool(0)
}

func (m *MockBoard) disableLocation(location Point) error {
    args := m.Called(location)
    return args.Error(0)
}

func (m *MockBoard) getVulnerables(color string) ([]Point, error) {
    args := m.Called(color)
    return args.Get(0).([]Point), args.Error(1)
}

func (m *MockBoard) setVulnerables(color string, vulnerables []Point) error {
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

func (m *MockBoard) possibleEnPassant(color string, location Point) ([]*EnPassant, error) {
    args := m.Called(color, location)
    return args.Get(0).([]*EnPassant), args.Error(1)
}

func (m *MockBoard) clearEnPassant(color string) error {
    args := m.Called(color)
    return args.Error(0)
}

func (m *MockBoard) ValidMoves(fromLocation Point) ([]Move, error) {
    args := m.Called(fromLocation)
    return args.Get(0).([]Move), args.Error(1)
}
 
func (m *MockBoard) Move(fromLocation Point, toLocation Point, promotion string) (Move, error) {
    args := m.Called(fromLocation, toLocation, promotion)
    return args.Get(0).(Move), args.Error(1)
}

func (m *MockBoard) CalculateMoves() error {
    args := m.Called()
    return args.Error(0)
}

func (m *MockBoard) CalculateMovesPartial(move Move) error {
    args := m.Called(move)
    return args.Error(0)
}

func (m *MockBoard) AvailableMoves(color string) ([]Move, error) {
    args := m.Called(color)
    return args.Get(0).([]Move), args.Error(1)
}

func (m *MockBoard) LegalMoves(color string) ([]Move, error) {
    args := m.Called(color)
    return args.Get(0).([]Move), args.Error(1)
}

func (m *MockBoard) Print() string {
    args := m.Called()
    return args.String(0)
}

func (m *MockBoard) Copy() (Board, error) {
    args := m.Called()
    return args.Get(0).(Board), args.Error(1)
}

func (m *MockBoard) UniqueString() string {
    args := m.Called()
    return args.String(0)
}

func (m *MockBoard) State() *BoardData {
    args := m.Called()
    return args.Get(0).(*BoardData)
}

func (m *MockBoard) Check(color string) bool {
    args := m.Called()
    return args.Bool(0)
}

func (m *MockBoard) Checkmate(color string) bool {
    args := m.Called(color)
    return args.Bool(0)
}

func (m *MockBoard) Stalemate(color string) bool {
    args := m.Called(color)
    return args.Bool(0)
}

func (m *MockBoard) Mate(color string) (bool, bool, error) {
    args := m.Called(color)
    return args.Bool(0), args.Bool(1), args.Error(2)
}

func (m *MockBoard) getPieceLocations() map[string][]Point {
    args := m.Called()
    return args.Get(0).(map[string][]Point)
}

func Test_NewSimpleBoard_DefaultFen(t *testing.T) {
    s, err := createSimpleBoardWithDefaultPieceLocations()
	assert.Nil(t, err)

	for y := 2; y <= 5; y++ {
		for x := 0; x <= 7; x++ {
			piece, ok := s.getPiece(Point{x, y})
			assert.True(t, ok)
			assert.Nil(t, piece)
		}
	}

	for _, y := range []int{1, 6} {
		for x := 0; x <= 7; x++ {
			piece, ok := s.getPiece(Point{x, y})
			assert.Nil(t, err)
			_, ok = piece.(*Pawn)
			assert.True(t, ok)
		}
	}

	for _, y := range []int{0, 7} {
		piece, ok := s.getPiece(Point{0, y})
		assert.True(t, ok)
		_, ok = piece.(*Rook)
		assert.True(t, ok)

		piece, ok = s.getPiece(Point{1, y})
		assert.True(t, ok)
		_, ok = piece.(*Knight)
		assert.True(t, ok)

		piece, ok = s.getPiece(Point{2, y})
		assert.True(t, ok)
		_, ok = piece.(*Bishop)
		assert.True(t, ok)

		piece, ok = s.getPiece(Point{3, y})
		assert.True(t, ok)
		_, ok = piece.(*Queen)
		assert.True(t, ok)

		piece, ok = s.getPiece(Point{4, y})
		assert.True(t, ok)
		_, ok = piece.(*King)
		assert.True(t, ok)

		piece, ok = s.getPiece(Point{5, y})
		assert.True(t, ok)
		_, ok = piece.(*Bishop)
		assert.True(t, ok)

		piece, ok = s.getPiece(Point{6, y})
		assert.True(t, ok)
		_, ok = piece.(*Knight)
		assert.True(t, ok)

		piece, ok = s.getPiece(Point{7, y})
		assert.True(t, ok)
		_, ok = piece.(*Rook)
		assert.True(t, ok)
	}
}

func Test_getAndSetPiece(t *testing.T) {
    s, err := newSimpleBoard(Point{8, 8})
    assert.Nil(t, err)

    p := newPawn("white", false, 0, 1)
    location := Point{0, 0}

    // adding piece to board
    ok := s.setPiece(location, p)
    assert.True(t, ok)

    piece, ok := s.getPiece(location)
    assert.True(t, ok)
    assert.Equal(t, p, piece)
    pieceLocations, ok := s.pieceLocationsMap["white"]
    assert.True(t, ok)
    exists := slices.Contains(pieceLocations, location)
    assert.True(t, exists)

    // removing piece from board
    ok = s.setPiece(location, nil)
    assert.True(t, ok)

    piece, ok = s.getPiece(location)
    assert.True(t, ok)
    assert.Nil(t, piece)
    pieceLocations, ok = s.pieceLocationsMap["white"]
    assert.True(t, ok)
    exists = slices.Contains(pieceLocations, location)
    assert.False(t, exists)
}

func Test_CalculateMoves_default(t *testing.T) {
    s, err := createSimpleBoardWithDefaultPieceLocations()
    assert.Nil(t, err)

    s.CalculateMoves()

    fromMoveCount := 0
    for _, toToMoveMap := range s.fromToToToMoveMap {
        fromMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 80, fromMoveCount)

    toMoveCount := 0
    for _, toToMoveMap := range s.toToFromToMoveMap {
        toMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 80, toMoveCount)
    
    whiteMoveKeys, err := s.AvailableMoves("white")
    assert.Nil(t, err)
    assert.Equal(t, 20, len(whiteMoveKeys))

    blackMoveKeys, err := s.AvailableMoves("black")
    assert.Nil(t, err)
    assert.Equal(t, 20, len(blackMoveKeys))

    assert.False(t, s.Check("white"))
    assert.False(t, s.Check("black"))
    assert.False(t, s.Checkmate("white"))
    assert.False(t, s.Checkmate("black"))
    assert.False(t, s.Stalemate("white"))
    assert.False(t, s.Stalemate("black"))
}

func Test_CalculateMoves_check(t *testing.T) {
    s, err := newSimpleBoard(Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(Point{0, 1}, newQueen("black"))
    s.setPiece(Point{0, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves()

    fromMoveCount := 0
    for _, toToMoveMap := range s.fromToToToMoveMap {
        fromMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 25, fromMoveCount)

    toMoveCount := 0
    for _, toToMoveMap := range s.toToFromToMoveMap {
        toMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 25, toMoveCount)
    
    whiteMoveKeys, err := s.AvailableMoves("white")
    assert.Nil(t, err)
    assert.Equal(t, 1, len(whiteMoveKeys))

    blackMoveKeys, err := s.AvailableMoves("black")
    assert.Nil(t, err)
    assert.Equal(t, 23, len(blackMoveKeys))

    assert.True(t, s.Check("white"))
    assert.False(t, s.Check("black"))
    assert.False(t, s.Checkmate("white"))
    assert.False(t, s.Checkmate("black"))
    assert.False(t, s.Stalemate("white"))
    assert.False(t, s.Stalemate("black"))
}

func Test_CalculateMoves_checkmate(t *testing.T) {
    s, err := newSimpleBoard(Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(Point{0, 1}, newQueen("black"))
    s.setPiece(Point{0, 2}, newKing("black", false, 0, -1))

    s.CalculateMoves()

    fromMoveCount := 0
    for _, toToMoveMap := range s.fromToToToMoveMap {
        fromMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 20, fromMoveCount)

    toMoveCount := 0
    for _, toToMoveMap := range s.toToFromToMoveMap {
        toMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 20, toMoveCount)
    
    whiteMoveKeys, err := s.AvailableMoves("white")
    assert.Nil(t, err)
    assert.Equal(t, 0, len(whiteMoveKeys))

    blackMoveKeys, err := s.AvailableMoves("black")
    assert.Nil(t, err)
    assert.Equal(t, 18, len(blackMoveKeys))

    assert.True(t, s.Check("white"))
    assert.False(t, s.Check("black"))
    assert.True(t, s.Checkmate("white"))
    assert.False(t, s.Checkmate("black"))
    assert.False(t, s.Stalemate("white"))
    assert.False(t, s.Stalemate("black"))
}

func Test_CalculateMoves_stalemate(t *testing.T) {
    s, err := newSimpleBoard(Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(Point{1, 2}, newQueen("black"))
    s.setPiece(Point{0, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves()

    fromMoveCount := 0
    for _, toToMoveMap := range s.fromToToToMoveMap {
        fromMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 26, fromMoveCount)

    toMoveCount := 0
    for _, toToMoveMap := range s.toToFromToMoveMap {
        toMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 26, toMoveCount)
    
    whiteMoveKeys, err := s.AvailableMoves("white")
    assert.Nil(t, err)
    assert.Equal(t, 0, len(whiteMoveKeys))

    blackMoveKeys, err := s.AvailableMoves("black")
    assert.Nil(t, err)
    assert.Equal(t, 26, len(blackMoveKeys))

    assert.False(t, s.Check("white"))
    assert.False(t, s.Check("black"))
    assert.False(t, s.Checkmate("white"))
    assert.False(t, s.Checkmate("black"))
    assert.True(t, s.Stalemate("white"))
    assert.False(t, s.Stalemate("black"))
}

func Test_CalculateMoves_noCastleThroughCheck(t *testing.T) {
    s, err := newSimpleBoard(Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(Point{4, 0}, newKing("white", false, 0, 1))
    s.setPiece(Point{0, 0}, newRook("white", false))
    s.setPiece(Point{3, 7}, newRook("black", false))
    s.setPiece(Point{4, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves()

    fromMoveCount := 0
    for _, toToMoveMap := range s.fromToToToMoveMap {
        fromMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 30, fromMoveCount)

    toMoveCount := 0
    for _, toToMoveMap := range s.toToFromToMoveMap {
        toMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 30, toMoveCount)
    
    whiteMoveKeys, err := s.AvailableMoves("white")
    assert.Nil(t, err)
    assert.Equal(t, 13, len(whiteMoveKeys))

    blackMoveKeys, err := s.AvailableMoves("black")
    assert.Nil(t, err)
    assert.Equal(t, 15, len(blackMoveKeys))

    assert.False(t, s.Check("white"))
    assert.False(t, s.Check("black"))
    assert.False(t, s.Checkmate("white"))
    assert.False(t, s.Checkmate("black"))
    assert.False(t, s.Stalemate("white"))
    assert.False(t, s.Stalemate("black"))
}

func Test_CalculateMoves_castle(t *testing.T) {
    s, err := newSimpleBoard(Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(Point{4, 0}, newKing("white", false, 0, 1))
    s.setPiece(Point{0, 0}, newRook("white", false))
    s.setPiece(Point{4, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves()

    fromMoveCount := 0
    for _, toToMoveMap := range s.fromToToToMoveMap {
        fromMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 22, fromMoveCount)

    toMoveCount := 0
    for _, toToMoveMap := range s.toToFromToMoveMap {
        toMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 22, toMoveCount)
    
    whiteMoveKeys, err := s.AvailableMoves("white")
    assert.Nil(t, err)
    assert.Equal(t, 16, len(whiteMoveKeys))

    blackMoveKeys, err := s.AvailableMoves("black")
    assert.Nil(t, err)
    assert.Equal(t, 5, len(blackMoveKeys))

    assert.False(t, s.Check("white"))
    assert.False(t, s.Check("black"))
    assert.False(t, s.Checkmate("white"))
    assert.False(t, s.Checkmate("black"))
    assert.False(t, s.Stalemate("white"))
    assert.False(t, s.Stalemate("black"))
}

func Test_CalculateMoves_promotion(t *testing.T) {
    s, err := newSimpleBoard(Point{8, 8})
    assert.Nil(t, err)

    s.setPiece(Point{7, 6}, newPawn("white", true, 0, 1))
    s.setPiece(Point{0, 0}, newKing("white", false, 0, 1))
    s.setPiece(Point{0, 7}, newKing("black", false, 0, -1))

    s.CalculateMoves()

    fromMoveCount := 0
    for _, toToMoveMap := range s.fromToToToMoveMap {
        fromMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 7, fromMoveCount)

    toMoveCount := 0
    for _, toToMoveMap := range s.toToFromToMoveMap {
        toMoveCount += len(toToMoveMap)
    }
    assert.Equal(t, 7, toMoveCount)
    
    whiteMoveKeys, err := s.AvailableMoves("white")
    assert.Nil(t, err)
    assert.Equal(t, 4, len(whiteMoveKeys))

    blackMoveKeys, err := s.AvailableMoves("black")
    assert.Nil(t, err)
    assert.Equal(t, 3, len(blackMoveKeys))

    assert.False(t, s.Check("white"))
    assert.False(t, s.Check("black"))
    assert.False(t, s.Checkmate("white"))
    assert.False(t, s.Checkmate("black"))
    assert.False(t, s.Stalemate("white"))
    assert.False(t, s.Stalemate("black"))
}

// TODO add en passant stuff to the unique string
// TODO use shortened color names for the unique string
func Test_MinimumString_Default(t *testing.T) {
    s, err := createSimpleBoardWithDefaultPieceLocations()
    assert.Nil(t, err)

    expected := "RblackmNblackBblackQblackKblackmBblackNblackRblackmPblackmPblackmPblackmPblackmPblackmPblackmPblackmPblackm32PwhitemPwhitemPwhitemPwhitemPwhitemPwhitemPwhitemPwhitemRwhitemNwhiteBwhiteQwhiteKwhitemBwhiteNwhiteRwhitem"
    actual := s.UniqueString()
    assert.Equal(t, expected, actual)
}

