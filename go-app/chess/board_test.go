package chess

import (
	"testing"
    //"slices"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBoard struct {
	mock.Mock
}

func (m *MockBoard) disablePieces(color int, disable bool) {
    m.Called(color, disable)
}

func (m *MockBoard) disableLocation(location *Point) {
    m.Called(location)
}

func (m *MockBoard) getPieceLocations() []Array100[Point] {
    args := m.Called()
    return args.Get(0).([]Array100[Point])
}

func (m *MockBoard) getPiece(location *Point) *Piece {
	args := m.Called(location)
    return args.Get(0).(*Piece)
}

func (m *MockBoard) getVulnerable(color int) *Vulnerable {
    args := m.Called(color)
    return args.Get(0).(*Vulnerable)
}

func (m *MockBoard) getEnPassant(color int) *EnPassant {
	args := m.Called(color)
    return args.Get(0).(*EnPassant)
}

func (m *MockBoard) possibleEnPassant(color int, location *Point) (*EnPassant, *EnPassant) {
    args := m.Called(color, location)
    return args.Get(0).(*EnPassant), args.Get(1).(*EnPassant)
}

func (m *MockBoard) MovesOfColor(color int) *Array100[FastMove] {
    args := m.Called(color)
    return args.Get(0).(*Array100[FastMove])
}

func (m *MockBoard) MovesOfLocation(fromLocation *Point) *Array100[FastMove] {
    args := m.Called(fromLocation)
    return args.Get(0).(*Array100[FastMove])
}

func (m *MockBoard) LegalMovesOfColor(color int) ([]FastMove, error) {
    args := m.Called(color)
    return args.Get(0).([]FastMove), args.Error(1)
}

func (m *MockBoard) LegalMovesOfLocation(fromLocation *Point) ([]FastMove, error) {
    args := m.Called(fromLocation)
    return args.Get(0).([]FastMove), args.Error(1)
}

func (m *MockBoard) CalculateMoves() {
    m.Called()
}

func (m *MockBoard) Check(color int) bool {
    args := m.Called(color)
    return args.Bool(0)
}

func (m *MockBoard) CheckmateAndStalemate(color int) (bool, bool, error) {
    args := m.Called(color)
    return args.Bool(0), args.Bool(1), args.Error(2)
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
			piece := b.getPiece(&Point{x, y})
			assert.True(t, piece != nil)
			assert.True(t, !piece.valid())
		}
	}

	for _, y := range []int{1, 6} {
		for x := 0; x <= 7; x++ {
			piece := b.getPiece(&Point{x, y})
            assert.True(t, piece != nil)
            assert.True(t, piece.index < 9) // is pawn
		}
	}

	for _, y := range []int{0, 7} {
		piece := b.getPiece(&Point{0, y})
		assert.True(t, piece != nil)
        assert.True(t, piece.index == 11 || piece.index == 12) // is rook

		piece = b.getPiece(&Point{1, y})
		assert.True(t, piece != nil)
        assert.True(t, piece.index == 9) // is knight

		piece = b.getPiece(&Point{2, y})
		assert.True(t, piece != nil)
        assert.True(t, piece.index == 10) // is bishop

		piece = b.getPiece(&Point{3, y})
		assert.True(t, piece != nil)
        assert.True(t, piece.index == 13) // is queen

		piece = b.getPiece(&Point{4, y})
		assert.True(t, piece != nil)
        assert.True(t, piece.index > 13) // is king

		piece = b.getPiece(&Point{5, y})
		assert.True(t, piece != nil)
        assert.True(t, piece.index == 10) // is bishop

		piece = b.getPiece(&Point{6, y})
		assert.True(t, piece != nil)
        assert.True(t, piece.index == 9) // is knight

		piece = b.getPiece(&Point{7, y})
		assert.True(t, piece != nil)
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
    setPiece(b, location, p)

    piece := b.getPiece(&location)
    assert.True(t, piece != nil)
    assert.Equal(t, p, *piece)

    //pieceLocations := b.pieceLocations[white]
    //exists := slices.Contains(pieceLocations, location)
    //assert.True(t, exists)

    // removing piece from board
    setPiece(b, location, Piece{0, 0})

    piece = b.getPiece(&location)
    assert.True(t, piece != nil)
    assert.True(t, !piece.valid())

    //pieceLocations = b.pieceLocations[white]
    //exists = slices.Contains(pieceLocations, location)
    //assert.False(t, exists)
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

    setPiece(b, Point{0, 0}, Piece{white, KING_D})
    setPiece(b, Point{0, 1}, Piece{black, QUEEN})
    setPiece(b, Point{0, 7}, Piece{black, KING_U})

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

    setPiece(b, Point{0, 0}, Piece{white, KING_D})
    setPiece(b, Point{0, 1}, Piece{black, QUEEN})
    setPiece(b, Point{0, 2}, Piece{black, KING_U})

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

    setPiece(b, Point{0, 0}, Piece{white, KING_D})
    setPiece(b, Point{1, 2}, Piece{black, QUEEN})
    setPiece(b, Point{0, 7}, Piece{black, KING_U})

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

    setPiece(b, Point{4, 0}, Piece{white, KING_D})
    setPiece(b, Point{0, 0}, Piece{white, ROOK})
    setPiece(b, Point{4, 7}, Piece{black, KING_U})
    setPiece(b, Point{3, 7}, Piece{black, ROOK})

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

    setPiece(b, Point{4, 0}, Piece{white, KING_D})
    setPiece(b, Point{0, 0}, Piece{white, ROOK})
    setPiece(b, Point{4, 7}, Piece{black, KING_U})

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

    setPiece(b, Point{7, 6}, Piece{white, PAWN_D_M})
    setPiece(b, Point{0, 0}, Piece{white, KING_D})
    setPiece(b, Point{0, 7}, Piece{black, KING_U})

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

