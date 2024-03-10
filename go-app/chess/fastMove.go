package chess

import (
    "fmt"
)

func createMoveSimple(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
    newPiece Piece,
) (FastMove, error) {
	piece, ok := b.getPiece(fromLocation)
	if !ok {
		return FastMove{}, fmt.Errorf("no piece at fromLocation")
	}

	capturedPiece, ok := b.getPiece(toLocation)
	if !ok {
		return FastMove{}, fmt.Errorf("no piece at toLocation")
	}

	enPassant, err := b.getEnPassant(color)
	if err != nil {
		return FastMove{}, err
	}

    vulnerable, err := b.getVulnerable(color)
    if err != nil {
        return FastMove{}, err
    }

    if newPiece == nil {
        newPiece, err = piece.copy()
        if err != nil {
            return FastMove{}, err
        }
    }

    newPieceLocations := Array4[PieceLocation]{}
    newPieceLocations.append(PieceLocation{piece: nil, location: fromLocation})
    newPieceLocations.append(PieceLocation{piece: newPiece, location: toLocation})

    oldPieceLocations := Array4[PieceLocation]{}
    oldPieceLocations.append(PieceLocation{piece: piece, location: fromLocation})
    oldPieceLocations.append(PieceLocation{piece: capturedPiece, location: toLocation})

    return FastMove{
        b: b,
        fromLocation: fromLocation,
        toLocation: toLocation,
        color: color,
        allyDefense: false,

        newPieceLocation: newPieceLocations,
        oldPieceLocation: oldPieceLocations,

        newEnPassant: EnPassant{target: Point{-1, -1}, risk: Point{-1, -1}},
        oldEnPassant: enPassant,

        newVulnerable: Vulnerable{start: Point{-1, -1}, end: Point{-1, -1}},
        oldVulnerable: vulnerable,
    }, nil
}

func createMoveRevealEnPassant(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
    newPiece Piece,
    newEnPassant EnPassant,
) (FastMove, error) {
    m, err := createMoveSimple(b, color, fromLocation, toLocation, newPiece)
    if err != nil {
        return m, err
    }

    m.newEnPassant = newEnPassant

    return m, nil
}

func createMoveCaptureEnPassant(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
    newPiece Piece,
) (FastMove, error) {
    m, err := createMoveSimple(b, color, fromLocation, toLocation, newPiece)
    if err != nil {
        return m, err
    }

    possibleEnPassant, err := b.possibleEnPassant(color, toLocation)
    if err != nil {
        return m, err
    }

    if len(possibleEnPassant) > 2 {
        return m, fmt.Errorf("more than 2 possible en passant")
    }

    for _, enPassant := range possibleEnPassant {
        capturedPiece, ok := b.getPiece(enPassant.risk)
        if !ok {
            return m, fmt.Errorf("no piece at enPassant.pieceLocation")
        }

        m.newPieceLocation.append(PieceLocation{piece: nil, location: enPassant.risk})
        m.oldPieceLocation.append(PieceLocation{piece: capturedPiece, location: enPassant.risk})
    }

    return m, nil
}

func createMoveAllyDefense(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
) (FastMove, error) {
    return FastMove{
        b: b,
        fromLocation: fromLocation,
        toLocation: toLocation,
        color: color,
        allyDefense: true,

        newPieceLocation: Array4[PieceLocation]{},
        oldPieceLocation: Array4[PieceLocation]{},

        newEnPassant: EnPassant{target: Point{-1, -1}, risk: Point{-1, -1}},
        oldEnPassant: EnPassant{target: Point{-1, -1}, risk: Point{-1, -1}},

        newVulnerable: Vulnerable{start: Point{-1, -1}, end: Point{-1, -1}},
        oldVulnerable: Vulnerable{start: Point{-1, -1}, end: Point{-1, -1}},
    }, nil
}

func createMoveCastle(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
    toKingLocation Point,
    toRookLocation Point,
    newVulnerable Vulnerable,
) (FastMove, error) {
	king, ok := b.getPiece(fromLocation)
	if !ok {
		return FastMove{}, fmt.Errorf("no piece at fromLocation")
	}

	rook, ok := b.getPiece(toLocation)
	if !ok {
		return FastMove{}, fmt.Errorf("no piece at toLocation")
	}

	enPassant, err := b.getEnPassant(color)
	if err != nil {
		return FastMove{}, err
	}

    vulnerable, err := b.getVulnerable(color)
    if err != nil {
        return FastMove{}, err
    }

	newKing, err := king.copy()
    if err != nil {
        return FastMove{}, err
    }

	newRook, err := rook.copy()
    if err != nil {
        return FastMove{}, err
    }

    newPieceLocations := Array4[PieceLocation]{}
    newPieceLocations.append(PieceLocation{piece: nil, location: fromLocation})
    newPieceLocations.append(PieceLocation{piece: nil, location: toLocation})
    newPieceLocations.append(PieceLocation{piece: newKing, location: toKingLocation})
    newPieceLocations.append(PieceLocation{piece: newRook, location: toRookLocation})

    oldPieceLocations := Array4[PieceLocation]{}
    oldPieceLocations.append(PieceLocation{piece: nil, location: toKingLocation})
    oldPieceLocations.append(PieceLocation{piece: nil, location: toRookLocation})
    oldPieceLocations.append(PieceLocation{piece: king, location: fromLocation})
    oldPieceLocations.append(PieceLocation{piece: rook, location: toLocation})

    return FastMove{
        b: b,
        fromLocation: fromLocation,
        toLocation: toLocation,
        color: color,
        allyDefense: false,

        newPieceLocation: newPieceLocations,
        oldPieceLocation: oldPieceLocations,

        newEnPassant: EnPassant{target: Point{-1, -1}, risk: Point{-1, -1}},
        oldEnPassant: enPassant,

        newVulnerable: newVulnerable,
        oldVulnerable: vulnerable,
    }, nil
}

type FastMove struct {
    b Board
    fromLocation Point
    toLocation Point
    color int
    allyDefense bool

    newPieceLocation Array4[PieceLocation]
    oldPieceLocation Array4[PieceLocation]

    newEnPassant EnPassant
    oldEnPassant EnPassant

    newVulnerable Vulnerable
    oldVulnerable Vulnerable
}

func (m *FastMove) execute() error {
    for i := 0; i < m.newPieceLocation.count; i++ {
        pieceLocation := m.newPieceLocation.array[i]
        ok := m.b.setPiece(pieceLocation.location, pieceLocation.piece)
        if !ok {
            return fmt.Errorf("could not set piece")
        }
    }

    err := m.b.setEnPassant(m.color, m.newEnPassant)
    if err != nil {
        return err
    }

    err = m.b.setVulnerable(m.color, m.newVulnerable)
    if err != nil {
        return err
    }

    return nil
}

func (m *FastMove) undo() error {
    for i := 0; i < m.oldPieceLocation.count; i++ {
        pieceLocation := m.oldPieceLocation.array[i]
        ok := m.b.setPiece(pieceLocation.location, pieceLocation.piece)
        if !ok {
            return fmt.Errorf("could not set piece")
        }
    }

    err := m.b.setEnPassant(m.color, m.oldEnPassant)
    if err != nil {
        return err
    }

    err = m.b.setVulnerable(m.color, m.oldVulnerable)
    if err != nil {
        return err
    }

    return nil
}

