package chess

import (
    "fmt"
)

func createMoveSimple(
    b Board,
    fromPiece Piece,
    fromLocation Point,
    toPiece Piece,
    toLocation Point,
    newPiece Piece,
) FastMove {
	enPassant, err := b.getEnPassant(fromPiece.color)
	if err != nil {
        panic(err)
	}

    vulnerable, err := b.getVulnerable(fromPiece.color)
    if err != nil {
        panic(err)
    }

    var promotion string
    if newPiece.valid() {
        promotion = newPiece.print()
    } else {
        newPiece = fromPiece.copy()
        promotion = ""
    }

    newPieceLocations := Array4[PieceLocation]{}
    newPieceLocations.append(PieceLocation{piece: Piece{0, 0}, location: fromLocation})
    newPieceLocations.append(PieceLocation{piece: newPiece, location: toLocation})

    oldPieceLocations := Array4[PieceLocation]{}
    oldPieceLocations.append(PieceLocation{piece: fromPiece, location: fromLocation})
    oldPieceLocations.append(PieceLocation{piece: toPiece, location: toLocation})

    return FastMove{
        b: b,
        fromLocation: fromLocation,
        toLocation: toLocation,
        color: fromPiece.color,
        allyDefense: false,
        promotion: promotion,

        newPieceLocation: newPieceLocations,
        oldPieceLocation: oldPieceLocations,

        newEnPassant: EnPassant{target: Point{-1, -1}, risk: Point{-1, -1}},
        oldEnPassant: enPassant,

        newVulnerable: Vulnerable{start: Point{-1, -1}, end: Point{-1, -1}},
        oldVulnerable: vulnerable,
    }
}

func createMoveRevealEnPassant(
    b Board,
    fromPiece Piece,
    fromLocation Point,
    toPiece Piece,
    toLocation Point,
    newPiece Piece,
    newEnPassant EnPassant,
) FastMove {
    m := createMoveSimple(b, fromPiece, fromLocation, toPiece, toLocation, newPiece)

    m.newEnPassant = newEnPassant

    return m
}

func createMoveCaptureEnPassant(
    b Board,
    fromPiece Piece,
    fromLocation Point,
    toPiece Piece,
    toLocation Point,
    newPiece Piece,
    enPassants []EnPassant,
) FastMove {
    m := createMoveSimple(b, fromPiece, fromLocation, toPiece, toLocation, newPiece)

    if len(enPassants) > 2 {
        panic("more than 2 possible en passant")
    }

    for _, enPassant := range enPassants {
        capturedPiece, ok := b.getPiece(enPassant.risk)
        if !ok {
            panic("no piece at enPassant.pieceLocation")
        }

        m.newPieceLocation.append(PieceLocation{piece: Piece{0, 0}, location: enPassant.risk})
        m.oldPieceLocation.append(PieceLocation{piece: capturedPiece, location: enPassant.risk})
    }

    return m
}

func createMoveAllyDefense(
    b Board,
    fromPiece Piece,
    fromLocation Point,
    toLocation Point,
) FastMove {
    return FastMove{
        b: b,
        fromLocation: fromLocation,
        toLocation: toLocation,
        color: fromPiece.color,
        allyDefense: true,
        promotion: "",

        newPieceLocation: Array4[PieceLocation]{},
        oldPieceLocation: Array4[PieceLocation]{},

        newEnPassant: EnPassant{target: Point{-1, -1}, risk: Point{-1, -1}},
        oldEnPassant: EnPassant{target: Point{-1, -1}, risk: Point{-1, -1}},

        newVulnerable: Vulnerable{start: Point{-1, -1}, end: Point{-1, -1}},
        oldVulnerable: Vulnerable{start: Point{-1, -1}, end: Point{-1, -1}},
    }
}

func createMoveCastle(
    b Board,
    king Piece,
    fromLocation Point,
    rook Piece,
    toLocation Point,
    toKingLocation Point,
    toRookLocation Point,
    newVulnerable Vulnerable,
) FastMove {
	enPassant, err := b.getEnPassant(king.color)
	if err != nil {
        panic(err)
	}

    vulnerable, err := b.getVulnerable(king.color)
    if err != nil {
        panic(err)
    }

	newKing := king.copy()
	newRook := rook.copy()

    newPieceLocations := Array4[PieceLocation]{}
    newPieceLocations.append(PieceLocation{piece: Piece{0, 0}, location: fromLocation})
    newPieceLocations.append(PieceLocation{piece: Piece{0, 0}, location: toLocation})
    newPieceLocations.append(PieceLocation{piece: newKing, location: toKingLocation})
    newPieceLocations.append(PieceLocation{piece: newRook, location: toRookLocation})

    oldPieceLocations := Array4[PieceLocation]{}
    oldPieceLocations.append(PieceLocation{piece: Piece{0, 0}, location: toKingLocation})
    oldPieceLocations.append(PieceLocation{piece: Piece{0, 0}, location: toRookLocation})
    oldPieceLocations.append(PieceLocation{piece: king, location: fromLocation})
    oldPieceLocations.append(PieceLocation{piece: rook, location: toLocation})

    return FastMove{
        b: b,
        fromLocation: fromLocation,
        toLocation: toLocation,
        color: king.color,
        allyDefense: false,
        promotion: "",

        newPieceLocation: newPieceLocations,
        oldPieceLocation: oldPieceLocations,

        newEnPassant: EnPassant{target: Point{-1, -1}, risk: Point{-1, -1}},
        oldEnPassant: enPassant,

        newVulnerable: newVulnerable,
        oldVulnerable: vulnerable,
    }
}

type FastMove struct {
    b Board
    fromLocation Point
    toLocation Point
    color int
    allyDefense bool
    promotion string

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

