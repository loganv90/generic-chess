package chess

import (
    "fmt"
)

func createSimpleMove(
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

	enPassantTarget, enPassantRisk, err := b.getEnPassant2(color)
	if err != nil {
		return FastMove{}, err
	}

    vulnerableStart, vulnerableEnd, err := b.getVulnerables2(color)
    if err != nil {
        return FastMove{}, err
    }

    if newPiece == nil {
        newPiece, err = piece.copy()
        if err != nil {
            return FastMove{}, err
        }
    }

    newPieceAndLocations := Array4[pieceAndLocation]{}
    newPieceAndLocations.append(pieceAndLocation{piece: nil, location: fromLocation})
    newPieceAndLocations.append(pieceAndLocation{piece: newPiece, location: toLocation})

    oldPieceAndLocations := Array4[pieceAndLocation]{}
    oldPieceAndLocations.append(pieceAndLocation{piece: piece, location: fromLocation})
    oldPieceAndLocations.append(pieceAndLocation{piece: capturedPiece, location: toLocation})

    return FastMove{
        b: b,
        fromLocation: fromLocation,
        toLocation: toLocation,
        color: color,
        allyDefense: false,

        newPieceAndLocation: newPieceAndLocations,
        oldPieceAndLocation: oldPieceAndLocations,

        newEnPassantTarget: Point{-1, -1},
        oldEnPassantTarget: enPassantTarget,

        newEnPassantRisk: Point{-1, -1},
        oldEnPassantRisk: enPassantRisk,

        newVulnerableStart: Point{-1, -1},
        oldVulnerableStart: vulnerableStart,

        newVulnerableEnd: Point{-1, -1},
        oldVulnerableEnd: vulnerableEnd,
    }, nil
}

func createRevealEnPassantMove(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
    target Point,
    newPiece Piece,
) (FastMove, error) {
    m, err := createSimpleMove(b, color, fromLocation, toLocation, newPiece)
    if err != nil {
        return m, err
    }

    m.newEnPassantTarget = target
    m.oldEnPassantRisk = toLocation

    return m, nil
}

func createCaptureEnPassantMove(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
    newPiece Piece,
) (FastMove, error) {
    m, err := createSimpleMove(b, color, fromLocation, toLocation, newPiece)
    if err != nil {
        return m, err
    }

    // TODO fix
    possibleEnPassant, err := b.possibleEnPassant(color, toLocation)
    if err == nil {
        for _, enPassant := range possibleEnPassant {
            capturedPiece, ok := b.getPiece(enPassant.pieceLocation)
            if !ok {
                return m, fmt.Errorf("no piece at enPassant.pieceLocation")
            }

            m.newPieceAndLocation.append(pieceAndLocation{piece: nil, location: enPassant.pieceLocation})
            m.oldPieceAndLocation.append(pieceAndLocation{piece: capturedPiece, location: enPassant.pieceLocation})
        }
    }

    return m, nil
}

func createAllyDefenseMove(
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
    }, nil
}

func createCastleMove(
    b Board,
    color int,
    fromLocation Point,
    toLocation Point,
    toKingLocation Point,
    toRookLocation Point,
    newVulnerableStart Point,
    newVulnerableEnd Point,
) (FastMove, error) {
	king, ok := b.getPiece(fromLocation)
	if !ok {
		return FastMove{}, fmt.Errorf("no piece at fromLocation")
	}

	rook, ok := b.getPiece(toLocation)
	if !ok {
		return FastMove{}, fmt.Errorf("no piece at toLocation")
	}

	enPassantTarget, enPassantRisk, err := b.getEnPassant2(color)
	if err != nil {
		return FastMove{}, err
	}

    vulnerableStart, vulnerableEnd, err := b.getVulnerables2(color)
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

    newPieceAndLocations := Array4[pieceAndLocation]{}
    newPieceAndLocations.append(pieceAndLocation{piece: newKing, location: toKingLocation})
    newPieceAndLocations.append(pieceAndLocation{piece: newRook, location: toRookLocation})

    oldPieceAndLocations := Array4[pieceAndLocation]{}
    oldPieceAndLocations.append(pieceAndLocation{piece: king, location: fromLocation})
    oldPieceAndLocations.append(pieceAndLocation{piece: rook, location: toLocation})

    return FastMove{
        b: b,
        fromLocation: fromLocation,
        toLocation: toLocation,
        color: color,
        allyDefense: false,
        newPieceAndLocation: newPieceAndLocations,
        oldPieceAndLocation: oldPieceAndLocations,
        newEnPassantTarget: Point{-1, -1},
        oldEnPassantTarget: enPassantTarget,
        newEnPassantRisk: Point{-1, -1},
        oldEnPassantRisk: enPassantRisk,
        newVulnerableStart: newVulnerableStart,
        oldVulnerableStart: vulnerableStart,
        newVulnerableEnd: newVulnerableEnd,
        oldVulnerableEnd: vulnerableEnd,
    }, nil
}

type FastMove struct {
    b Board
    fromLocation Point
    toLocation Point
    color int
    allyDefense bool

    newPieceAndLocation Array4[pieceAndLocation]
    oldPieceAndLocation Array4[pieceAndLocation]

    newEnPassantTarget Point
    oldEnPassantTarget Point

    newEnPassantRisk Point
    oldEnPassantRisk Point

    newVulnerableStart Point
    oldVulnerableStart Point

    newVulnerableEnd Point
    oldVulnerableEnd Point
}

func (m *FastMove) execute() error {
    for i := 0; i < m.newPieceAndLocation.count; i++ {
        pieceAndLocation := m.newPieceAndLocation.array[i]
        ok := m.b.setPiece(pieceAndLocation.location, pieceAndLocation.piece)
        if !ok {
            return fmt.Errorf("could not set piece")
        }
    }

    err := m.b.setEnPassant2(m.color, m.newEnPassantTarget, m.newEnPassantRisk)
    if err != nil {
        return err
    }

    err = m.b.setVulnerables2(m.color, m.newVulnerableStart, m.newVulnerableEnd)
    if err != nil {
        return err
    }

    return nil
}

func (m *FastMove) undo() error {
    for i := 0; i < m.oldPieceAndLocation.count; i++ {
        pieceAndLocation := m.oldPieceAndLocation.array[i]
        ok := m.b.setPiece(pieceAndLocation.location, pieceAndLocation.piece)
        if !ok {
            return fmt.Errorf("could not set piece")
        }
    }

    err := m.b.setEnPassant2(m.color, m.oldEnPassantTarget, m.oldEnPassantRisk)
    if err != nil {
        return err
    }

    err = m.b.setVulnerables2(m.color, m.oldVulnerableStart, m.oldVulnerableEnd)
    if err != nil {
        return err
    }

    return nil
}

