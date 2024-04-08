package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pawn_Moves_Unmoved(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(7, 7, 2)
    assert.Nil(t, err)

    pawn := b.getAllPiece(white, PAWN_D)
    b.setPiece(b.getIndex(3, 3), pawn)

	pawn.moves(b, b.getIndex(3, 3))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 2,
        &b.captureMoves[white], 0,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(3, 4),
            b.getIndex(3, 5),
        },
    )
}

func Test_Pawn_Moves_Moved(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(7, 7, 2)
    assert.Nil(t, err)

    pawn := b.getAllPiece(white, PAWN_D_M)
    b.setPiece(b.getIndex(3, 3), pawn)

	pawn.moves(b, b.getIndex(3, 3))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 1,
        &b.captureMoves[white], 0,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(3, 4),
        },
    )
}

func Test_Pawn_Moves_Capturing(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(7, 7, 2)
    assert.Nil(t, err)

    pawn := b.getAllPiece(white, PAWN_D)
    blackPawn := b.getAllPiece(black, PAWN_D)
    b.setPiece(b.getIndex(3, 3), pawn)
    b.setPiece(b.getIndex(4, 4), blackPawn)

	pawn.moves(b, b.getIndex(3, 3))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 2,
        &b.captureMoves[white], 1,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(3, 4),
            b.getIndex(3, 5),
            b.getIndex(4, 4),
        },
    )
}

func Test_Pawn_Moves_CapturingEnPassant(t *testing.T) {
    white := 0
    black := 1

    b, err := newSimpleBoard(7, 7, 2)
    assert.Nil(t, err)

    pawn := b.getAllPiece(white, PAWN_D)
    b.setPiece(b.getIndex(3, 3), pawn)
    b.setEnPassant(black, b.getIndex(4, 4), b.getIndex(3, 4))

    pawn.moves(b, b.getIndex(3, 3))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 2,
        &b.captureMoves[white], 1,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(3, 4),
            b.getIndex(3, 5),
            b.getIndex(4, 4),
        },
    )
}

func Test_Pawn_Moves_Promotion(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(5, 5, 2)
    assert.Nil(t, err)

    pawn := b.getAllPiece(white, PAWN_D)
    b.setPiece(b.getIndex(2, 2), pawn)

	pawn.moves(b, b.getIndex(2, 2))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 5,
        &b.captureMoves[white], 0,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(2, 3),
            b.getIndex(2, 4),
        },
    )
}

func Test_Knight_Moves(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(5, 5, 2)
    assert.Nil(t, err)

    knight := b.getAllPiece(white, KNIGHT)
    b.setPiece(b.getIndex(2, 2), knight)

    knight.moves(b, b.getIndex(2, 2))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 8,
        &b.captureMoves[white], 0,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(0, 1),
            b.getIndex(0, 3),
            b.getIndex(1, 0),
            b.getIndex(1, 4),
            b.getIndex(3, 0),
            b.getIndex(3, 4),
            b.getIndex(4, 1),
            b.getIndex(4, 3),
        },
    )
}

func Test_Bishop_Moves(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(5, 5, 2)
    assert.Nil(t, err)

    bishop := b.getAllPiece(white, BISHOP)
    b.setPiece(b.getIndex(2, 2), bishop)

	bishop.moves(b, b.getIndex(2, 2))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 8,
        &b.captureMoves[white], 0,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(0, 0),
            b.getIndex(1, 1),
            b.getIndex(3, 3),
            b.getIndex(4, 4),
            b.getIndex(0, 4),
            b.getIndex(1, 3),
            b.getIndex(3, 1),
            b.getIndex(4, 0),
        },
    )
}

func Test_Rook_Moves(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(5, 5, 2)
    assert.Nil(t, err)

    rook := b.getAllPiece(white, ROOK)
    b.setPiece(b.getIndex(2, 2), rook)

	rook.moves(b, b.getIndex(2, 2))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 8,
        &b.captureMoves[white], 0,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(2, 0),
            b.getIndex(2, 1),
            b.getIndex(2, 3),
            b.getIndex(2, 4),
            b.getIndex(0, 2),
            b.getIndex(1, 2),
            b.getIndex(3, 2),
            b.getIndex(4, 2),
        },
    )
}

func Test_Queen_Moves(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(5, 5, 2)
    assert.Nil(t, err)

    queen := b.getAllPiece(white, QUEEN)
    b.setPiece(b.getIndex(2, 2), queen)

    queen.moves(b, b.getIndex(2, 2))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 16,
        &b.captureMoves[white], 0,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(0, 0),
            b.getIndex(1, 1),
            b.getIndex(3, 3),
            b.getIndex(4, 4),
            b.getIndex(0, 4),
            b.getIndex(1, 3),
            b.getIndex(3, 1),
            b.getIndex(4, 0),
            b.getIndex(2, 0),
            b.getIndex(2, 1),
            b.getIndex(2, 3),
            b.getIndex(2, 4),
            b.getIndex(0, 2),
            b.getIndex(1, 2),
            b.getIndex(3, 2),
            b.getIndex(4, 2),
        },
    )
}

func Test_King_Moves_CanCastleAndUnmoved(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(5, 5, 2)
    assert.Nil(t, err)

    king := b.getAllPiece(white, KING_D)
    rook := b.getAllPiece(white, ROOK)
    b.setPiece(b.getIndex(2, 2), king)
    b.setPiece(b.getIndex(0, 2), rook)
    b.setPiece(b.getIndex(4, 2), rook)

    king.moves(b, b.getIndex(2, 2))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 10,
        &b.captureMoves[white], 0,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(2, 1),
            b.getIndex(2, 3),
            b.getIndex(0, 2),
            b.getIndex(1, 2),
            b.getIndex(3, 2),
            b.getIndex(4, 2),
            b.getIndex(1, 1),
            b.getIndex(3, 3),
            b.getIndex(3, 1),
            b.getIndex(1, 3),
        },
    )
}

func Test_King_Moves_CanCastleAndMoved(t *testing.T) {
    white := 0

    b, err := newSimpleBoard(5, 5, 2)
    assert.Nil(t, err)

    king := b.getAllPiece(white, KING_D_M)
    rook := b.getAllPiece(white, ROOK)
    b.setPiece(b.getIndex(2, 2), king)
    b.setPiece(b.getIndex(0, 2), rook)
    b.setPiece(b.getIndex(4, 2), rook)

    king.moves(b, b.getIndex(2, 2))

    Assert_LengthAndToLocations(
        t,
        &b.moves[white], 8,
        &b.captureMoves[white], 0,
        &b.defenseMoves[white], 0,
        []*Point{
            b.getIndex(2, 1),
            b.getIndex(2, 3),
            b.getIndex(1, 2),
            b.getIndex(3, 2),
            b.getIndex(1, 1),
            b.getIndex(3, 3),
            b.getIndex(3, 1),
            b.getIndex(1, 3),
        },
    )
}

func Assert_LengthAndToLocations(
    t *testing.T,
    moves *Array1000[FastMove],
    movesLength int,
    captureMoves *Array1000[FastMove],
    captureMovesLength int,
    defenseMoves *Array1000[FastMove],
    defenseMovesLength int,
    toLocations []*Point,
) {
    assert.Equal(t, movesLength, moves.count)
    assert.Equal(t, captureMovesLength, captureMoves.count)
    assert.Equal(t, defenseMovesLength, defenseMoves.count)

    moveMap := map[*Point]bool{}
    for i := 0; i < moves.count; i++ {
        m := moves.array[i]
        moveMap[m.toLocation] = true
    }
    for i := 0; i < captureMoves.count; i++ {
        m := captureMoves.array[i]
        moveMap[m.toLocation] = true
    }
    for i := 0; i < defenseMoves.count; i++ {
        m := defenseMoves.array[i]
        moveMap[m.toLocation] = true
    }

    for _, toLocation := range toLocations {
        assert.True(t, moveMap[toLocation])
    }
}

