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

func (m *MockBoard) getPiece(location Point) (Piece, bool) {
	args := m.Called(location)

	if args.Get(0) == nil {
		return Piece{0, 0}, args.Bool(1)
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

func (m *MockBoard) getVulnerable(color int) (Vulnerable, error) {
    args := m.Called(color)
    return args.Get(0).(Vulnerable), args.Error(1)
}

func (m *MockBoard) setVulnerable(color int, vulnerable Vulnerable) error {
    args := m.Called(color, vulnerable)
    return args.Error(0)
}

func (m *MockBoard) getEnPassant(color int) (EnPassant, error) {
	args := m.Called(color)
	return args.Get(0).(EnPassant), args.Error(1)
}

func (m *MockBoard) setEnPassant(color int, enPassant EnPassant) error {
    args := m.Called(color, enPassant)
    return args.Error(0)
}

func (m *MockBoard) possibleEnPassant(color int, location Point) ([]EnPassant, error) {
    args := m.Called(color, location)
    return args.Get(0).([]EnPassant), args.Error(1)
}

func (m *MockBoard) disablePieces(color int, disable bool) error {
    args := m.Called(color, disable)
    return args.Error(0)
}

func (m *MockBoard) getPieceLocations() [][]Point {
    args := m.Called()
    return args.Get(0).([][]Point)
}

func (m *MockBoard) CalculateMoves() error {
    args := m.Called()
    return args.Error(0)
}

func (m *MockBoard) MovesOfColor(color int) (*Array100[FastMove], error) {
    args := m.Called(color)
    return args.Get(0).(*Array100[FastMove]), args.Error(1)
}

func (m *MockBoard) MovesOfLocation(fromLocation Point) (*Array100[FastMove], error) {
    args := m.Called(fromLocation)
    return args.Get(0).(*Array100[FastMove]), args.Error(1)
}

func (m *MockBoard) LegalMovesOfColor(color int) ([]FastMove, error) {
    args := m.Called(color)
    return args.Get(0).([]FastMove), args.Error(1)
}

func (m *MockBoard) LegalMovesOfLocation(location Point) ([]FastMove, error) {
    args := m.Called(location)
    return args.Get(0).([]FastMove), args.Error(1)
}

func (m *MockBoard) CheckmateAndStalemate(color int) (bool, bool, error) {
    args := m.Called(color)
    return args.Bool(0), args.Bool(1), args.Error(2)
}

func (m *MockBoard) Check(color int) bool {
    args := m.Called(color)
    return args.Bool(0)
}

func (m *MockBoard) Print() string {
    args := m.Called()
    return args.String(0)
}

func (m *MockBoard) State() *BoardData {
    args := m.Called()
    return args.Get(0).(*BoardData)
}

func (m *MockBoard) Copy() (Board, error) {
    args := m.Called()
    return args.Get(0).(Board), args.Error(1)
}

func (m *MockBoard) UniqueString() string {
    args := m.Called()
    return args.String(0)
}

func Test_NewSimpleBoard_DefaultFen(t *testing.T) {
    b, err := createSimpleBoardWithDefaultPieceLocations()
	assert.Nil(t, err)

	for y := 2; y <= 5; y++ {
		for x := 0; x <= 7; x++ {
			piece, ok := b.getPiece(Point{x, y})
			assert.True(t, ok)
			assert.True(t, !piece.valid())
		}
	}

	for _, y := range []int{1, 6} {
		for x := 0; x <= 7; x++ {
			piece, ok := b.getPiece(Point{x, y})
            assert.True(t, ok)
            assert.True(t, piece.index < 9) // is pawn
		}
	}

	for _, y := range []int{0, 7} {
		piece, ok := b.getPiece(Point{0, y})
		assert.True(t, ok)
        assert.True(t, piece.index == 11 || piece.index == 12) // is rook

		piece, ok = b.getPiece(Point{1, y})
		assert.True(t, ok)
        assert.True(t, piece.index == 9) // is knight

		piece, ok = b.getPiece(Point{2, y})
		assert.True(t, ok)
        assert.True(t, piece.index == 10) // is bishop

		piece, ok = b.getPiece(Point{3, y})
		assert.True(t, ok)
        assert.True(t, piece.index == 13) // is queen

		piece, ok = b.getPiece(Point{4, y})
		assert.True(t, ok)
        assert.True(t, piece.index > 13) // is king

		piece, ok = b.getPiece(Point{5, y})
		assert.True(t, ok)
        assert.True(t, piece.index == 10) // is bishop

		piece, ok = b.getPiece(Point{6, y})
		assert.True(t, ok)
        assert.True(t, piece.index == 9) // is knight

		piece, ok = b.getPiece(Point{7, y})
		assert.True(t, ok)
        assert.True(t, piece.index == 11 || piece.index == 12) // is rook
	}
}

func Test_getAndSetPiece(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(Point{8, 8}, 2)
    assert.Nil(t, err)

    p := Piece{white, PAWN_D}
    location := Point{0, 0}

    // adding piece to board
    ok := b.setPiece(location, p)
    assert.True(t, ok)

    piece, ok := b.getPiece(location)
    assert.True(t, ok)
    assert.Equal(t, p, piece)

    pieceLocations := b.pieceLocations[white]
    exists := slices.Contains(pieceLocations, location)
    assert.True(t, exists)

    // removing piece from board
    ok = b.setPiece(location, Piece{0, 0})
    assert.True(t, ok)

    piece, ok = b.getPiece(location)
    assert.True(t, ok)
    assert.True(t, !piece.valid())

    pieceLocations = b.pieceLocations[white]
    exists = slices.Contains(pieceLocations, location)
    assert.False(t, exists)
}

func Test_CalculateMoves_default(t *testing.T) {
    white := 0
    black := 1

    b, err := createSimpleBoardWithDefaultPieceLocations()
    assert.Nil(t, err)

    b.CalculateMoves()

    whiteMoveKeys, err := b.LegalMovesOfColor(white)
    assert.Nil(t, err)
    assert.Equal(t, 20, len(whiteMoveKeys))

    blackMoveKeys, err := b.LegalMovesOfColor(black)
    assert.Nil(t, err)
    assert.Equal(t, 20, len(blackMoveKeys))

    checkmate, stalemate, err := b.CheckmateAndStalemate(white)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)

    checkmate, stalemate, err = b.CheckmateAndStalemate(black)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)
}

func Test_CalculateMoves_check(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(Point{8, 8}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{0, 0}, Piece{white, KING_D})
    b.setPiece(Point{0, 1}, Piece{black, QUEEN})
    b.setPiece(Point{0, 7}, Piece{black, KING_U})

    b.CalculateMoves()

    whiteMoveKeys, err := b.LegalMovesOfColor(white)
    assert.Nil(t, err)
    assert.Equal(t, 1, len(whiteMoveKeys))

    blackMoveKeys, err := b.LegalMovesOfColor(black)
    assert.Nil(t, err)
    assert.Equal(t, 23, len(blackMoveKeys))

    checkmate, stalemate, err := b.CheckmateAndStalemate(white)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)

    checkmate, stalemate, err = b.CheckmateAndStalemate(black)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)
}

func Test_CalculateMoves_checkmate(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(Point{8, 8}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{0, 0}, Piece{white, KING_D})
    b.setPiece(Point{0, 1}, Piece{black, QUEEN})
    b.setPiece(Point{0, 2}, Piece{black, KING_U})

    b.CalculateMoves()

    whiteMoveKeys, err := b.LegalMovesOfColor(white)
    assert.Nil(t, err)
    assert.Equal(t, 0, len(whiteMoveKeys))

    blackMoveKeys, err := b.LegalMovesOfColor(black)
    assert.Nil(t, err)
    assert.Equal(t, 18, len(blackMoveKeys))

    checkmate, stalemate, err := b.CheckmateAndStalemate(white)
    assert.Nil(t, err)
    assert.True(t, checkmate)
    assert.False(t, stalemate)

    checkmate, stalemate, err = b.CheckmateAndStalemate(black)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)
}

func Test_CalculateMoves_stalemate(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(Point{8, 8}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{0, 0}, Piece{white, KING_D})
    b.setPiece(Point{1, 2}, Piece{black, QUEEN})
    b.setPiece(Point{0, 7}, Piece{black, KING_U})

    b.CalculateMoves()

    whiteMoveKeys, err := b.LegalMovesOfColor(white)
    assert.Nil(t, err)
    assert.Equal(t, 0, len(whiteMoveKeys))

    blackMoveKeys, err := b.LegalMovesOfColor(black)
    assert.Nil(t, err)
    assert.Equal(t, 26, len(blackMoveKeys))

    checkmate, stalemate, err := b.CheckmateAndStalemate(white)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.True(t, stalemate)

    checkmate, stalemate, err = b.CheckmateAndStalemate(black)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)
}

func Test_CalculateMoves_noCastleThroughCheck(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(Point{8, 8}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{4, 0}, Piece{white, KING_D})
    b.setPiece(Point{0, 0}, Piece{white, ROOK})
    b.setPiece(Point{4, 7}, Piece{black, KING_U})
    b.setPiece(Point{3, 7}, Piece{black, ROOK})

    b.CalculateMoves()

    whiteMoveKeys, err := b.LegalMovesOfColor(white)
    assert.Nil(t, err)
    assert.Equal(t, 13, len(whiteMoveKeys))

    blackMoveKeys, err := b.LegalMovesOfColor(black)
    assert.Nil(t, err)
    assert.Equal(t, 15, len(blackMoveKeys))

    checkmate, stalemate, err := b.CheckmateAndStalemate(white)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)

    checkmate, stalemate, err = b.CheckmateAndStalemate(black)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)
}

func Test_CalculateMoves_castle(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(Point{8, 8}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{4, 0}, Piece{white, KING_D})
    b.setPiece(Point{0, 0}, Piece{white, ROOK})
    b.setPiece(Point{4, 7}, Piece{black, KING_U})

    b.CalculateMoves()

    whiteMoveKeys, err := b.LegalMovesOfColor(white)
    assert.Nil(t, err)
    assert.Equal(t, 16, len(whiteMoveKeys))

    blackMoveKeys, err := b.LegalMovesOfColor(black)
    assert.Nil(t, err)
    assert.Equal(t, 5, len(blackMoveKeys))

    checkmate, stalemate, err := b.CheckmateAndStalemate(white)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)

    checkmate, stalemate, err = b.CheckmateAndStalemate(black)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)
}

func Test_CalculateMoves_promotion(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(Point{8, 8}, 2)
    assert.Nil(t, err)

    b.setPiece(Point{7, 6}, Piece{white, PAWN_D_M})
    b.setPiece(Point{0, 0}, Piece{white, KING_D})
    b.setPiece(Point{0, 7}, Piece{black, KING_U})

    b.CalculateMoves()

    whiteMoveKeys, err := b.LegalMovesOfColor(white)
    assert.Nil(t, err)
    assert.Equal(t, 7, len(whiteMoveKeys))

    blackMoveKeys, err := b.LegalMovesOfColor(black)
    assert.Nil(t, err)
    assert.Equal(t, 3, len(blackMoveKeys))

    checkmate, stalemate, err := b.CheckmateAndStalemate(white)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)

    checkmate, stalemate, err = b.CheckmateAndStalemate(black)
    assert.Nil(t, err)
    assert.False(t, checkmate)
    assert.False(t, stalemate)
}

// TODO add en passant stuff to the unique string
func Test_MinimumString_Default(t *testing.T) {
    b, err := createSimpleBoardWithDefaultPieceLocations()
    assert.Nil(t, err)

    expected := "R1N1B1Q1K1B1N1R1P1P1P1P1P1P1P1P132P0P0P0P0P0P0P0P0R0N0B0Q0K0B0N0R0"
    actual := b.UniqueString()
    assert.Equal(t, expected, actual)
}

