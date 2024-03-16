package chess

import (
    "fmt"
)

func addMoveSimple(
    b Board,
    fromPiece *Piece,
    fromLocation *Point,
    toPiece *Piece,
    toLocation *Point,
    newPiece *Piece,
    moves *Array100[FastMove],
) {
	enPassant := b.getEnPassant(fromPiece.color)
	if enPassant == nil {
        panic("what")
	}

    vulnerable := b.getVulnerable(fromPiece.color)
    if vulnerable == nil {
        panic("what")
    }

    var promotion string
    if newPiece != nil {
        promotion = newPiece.print()
    } else {
        newPiece = fromPiece
        promotion = ""
    }


    newPieceLocations := Array4[PieceLocation]{}

    pieceLocation := newPieceLocations.get()
    pieceLocation.piece = Piece{0, 0}
    pieceLocation.location = *fromLocation
    newPieceLocations.next()

    pieceLocation = newPieceLocations.get()
    pieceLocation.piece = *newPiece
    pieceLocation.location = *toLocation
    newPieceLocations.next()


    oldPieceLocations := Array4[PieceLocation]{}

    pieceLocation = oldPieceLocations.get()
    pieceLocation.piece = *fromPiece
    pieceLocation.location = *fromLocation
    oldPieceLocations.next()

    pieceLocation = oldPieceLocations.get()
    pieceLocation.piece = *toPiece
    pieceLocation.location = *toLocation
    oldPieceLocations.next()


    move := moves.get()
    move.b = b
    move.fromLocation = *fromLocation
    move.toLocation = *toLocation
    move.color = fromPiece.color
    move.allyDefense = false
    move.promotion = promotion

    move.newPieceLocation = newPieceLocations
    move.oldPieceLocation = oldPieceLocations

    move.newEnPassant = EnPassant{target: Point{-1, -1}, risk: Point{-1, -1}}
    move.oldEnPassant = *enPassant

    move.newVulnerable = Vulnerable{start: Point{-1, -1}, end: Point{-1, -1}}
    move.oldVulnerable = *vulnerable
    moves.next()
}

func addMoveRevealEnPassant(
    b Board,
    fromPiece *Piece,
    fromLocation *Point,
    toPiece *Piece,
    toLocation *Point,
    newPiece *Piece,
    newEnPassant *EnPassant,
    moves *Array100[FastMove],
) {
    move := moves.get()

    addMoveSimple(b, fromPiece, fromLocation, toPiece, toLocation, newPiece, moves)

    move.newEnPassant = *newEnPassant
}

func addMoveCaptureEnPassant(
    b Board,
    fromPiece *Piece,
    fromLocation *Point,
    toPiece *Piece,
    toLocation *Point,
    newPiece *Piece,
    moves *Array100[FastMove],
    en1 *EnPassant,
    en2 *EnPassant,
) {
    move := moves.get()

    addMoveSimple(b, fromPiece, fromLocation, toPiece, toLocation, newPiece, moves)

    if en1 == nil {
        return
    }

    capturedPiece := b.getPiece(&en1.risk)

    if capturedPiece == nil {
        return
    }

    pieceLocation := move.newPieceLocation.get()
    pieceLocation.piece = Piece{0, 0}
    pieceLocation.location = en1.risk
    move.newPieceLocation.next()

    pieceLocation = move.oldPieceLocation.get()
    pieceLocation.piece = *capturedPiece
    pieceLocation.location = en1.risk
    move.oldPieceLocation.next()

    if en2 == nil {
        return
    }

    capturedPiece = b.getPiece(&en2.risk)

    if capturedPiece == nil {
        return
    }

    pieceLocation = move.newPieceLocation.get()
    pieceLocation.piece = Piece{0, 0}
    pieceLocation.location = en2.risk
    move.newPieceLocation.next()

    pieceLocation = move.oldPieceLocation.get()
    pieceLocation.piece = *capturedPiece
    pieceLocation.location = en2.risk
    move.oldPieceLocation.next()
}

func addMoveAllyDefense(
    b Board,
    fromPiece *Piece,
    fromLocation *Point,
    toLocation *Point,
    moves *Array100[FastMove],
) {
    move := moves.get()
    move.b = b
    move.fromLocation = *fromLocation
    move.toLocation = *toLocation
    move.color = fromPiece.color
    move.allyDefense = true
    move.promotion = ""
    moves.next()
}

func addMoveCastle(
    b Board,
    king *Piece,
    fromLocation *Point,
    toKingLocation *Point,
    rook *Piece,
    toLocation *Point,
    toRookLocation *Point,
    newVulnerable *Vulnerable,
    moves *Array100[FastMove],
) {
	enPassant := b.getEnPassant(king.color)
	if enPassant == nil {
        panic("what")
	}

    vulnerable := b.getVulnerable(king.color)
    if vulnerable == nil {
        panic("what")
    }

    newKing := Piece{0, 0}
    newRook := Piece{0, 0}
	king.copy(&newKing)
	rook.copy(&newRook)


    newPieceLocations := Array4[PieceLocation]{}

    pieceLocation := newPieceLocations.get()
    pieceLocation.piece = Piece{0, 0}
    pieceLocation.location = *fromLocation
    newPieceLocations.next()

    pieceLocation = newPieceLocations.get()
    pieceLocation.piece = Piece{0, 0}
    pieceLocation.location = *toLocation
    newPieceLocations.next()

    pieceLocation = newPieceLocations.get()
    pieceLocation.piece = newKing
    pieceLocation.location = *toKingLocation
    newPieceLocations.next()

    pieceLocation = newPieceLocations.get()
    pieceLocation.piece = newRook
    pieceLocation.location = *toRookLocation
    newPieceLocations.next()


    oldPieceLocations := Array4[PieceLocation]{}

    pieceLocation = oldPieceLocations.get()
    pieceLocation.piece = Piece{0, 0}
    pieceLocation.location = *toKingLocation
    oldPieceLocations.next()

    pieceLocation = oldPieceLocations.get()
    pieceLocation.piece = Piece{0, 0}
    pieceLocation.location = *toRookLocation
    oldPieceLocations.next()

    pieceLocation = oldPieceLocations.get()
    pieceLocation.piece = *king
    pieceLocation.location = *fromLocation
    oldPieceLocations.next()

    pieceLocation = oldPieceLocations.get()
    pieceLocation.piece = *rook
    pieceLocation.location = *toLocation
    oldPieceLocations.next()


    move := moves.get()
    move.b = b
    move.fromLocation = *fromLocation
    move.toLocation = *toLocation
    move.color = king.color
    move.allyDefense = false
    move.promotion = ""

    move.newPieceLocation = newPieceLocations
    move.oldPieceLocation = oldPieceLocations

    move.newEnPassant = EnPassant{target: Point{-1, -1}, risk: Point{-1, -1}}
    move.oldEnPassant = *enPassant

    move.newVulnerable = *newVulnerable
    move.oldVulnerable = *vulnerable
    moves.next()
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

        piece := m.b.getPiece(&pieceLocation.location)
        if piece == nil {
            return fmt.Errorf("could not get piece")
        }
        piece.color = pieceLocation.piece.color
        piece.index = pieceLocation.piece.index
    }

    enPassant := m.b.getEnPassant(m.color)
    if enPassant == nil {
        return fmt.Errorf("could not get en passant")
    }
    enPassant.target = m.newEnPassant.target
    enPassant.risk = m.newEnPassant.risk

    vulnerable := m.b.getVulnerable(m.color)
    if vulnerable == nil {
        return fmt.Errorf("could not get vulnerable")
    }
    vulnerable.start = m.newVulnerable.start
    vulnerable.end = m.newVulnerable.end

    return nil
}

func (m *FastMove) undo() error {
    for i := 0; i < m.oldPieceLocation.count; i++ {
        pieceLocation := m.oldPieceLocation.array[i]

        piece := m.b.getPiece(&pieceLocation.location)
        if piece == nil {
            return fmt.Errorf("could not get piece")
        }
        piece.color = pieceLocation.piece.color
        piece.index = pieceLocation.piece.index
    }

    enPassant := m.b.getEnPassant(m.color)
    if enPassant == nil {
        return fmt.Errorf("could not get en passant")
    }
    enPassant.target = m.oldEnPassant.target
    enPassant.risk = m.oldEnPassant.risk

    vulnerable := m.b.getVulnerable(m.color)
    if vulnerable == nil {
        return fmt.Errorf("could not get vulnerable")
    }
    vulnerable.start = m.oldVulnerable.start
    vulnerable.end = m.oldVulnerable.end

    return nil
}

