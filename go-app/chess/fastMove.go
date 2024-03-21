package chess

func addMoveSimple(
    b *SimpleBoard,
    fromPiece *Piece,
    fromLocation *Point,
    toPiece *Piece,
    toLocation *Point,
    newPiece *Piece,
    moves *Array1000[FastMove],
) {
    move := moves.get()

    move.b = b
    move.fromLocation = fromLocation
    move.toLocation = toLocation
    move.color = fromPiece.color
    move.allyDefense = false

	target, risk := b.getEnPassant(fromPiece.color)
    start, end := b.getVulnerable(fromPiece.color)

    if newPiece != nil { // promotion
        move.promotionIndex = newPiece.index
    } else { // no promotion
        move.promotionIndex = -1
        newPiece = b.movePiece(fromPiece)
    }

    move.newPiece.clear()
    move.newPiece.set(nil)
    move.newPiece.set(newPiece)

    move.oldPiece.clear()
    move.oldPiece.set(fromPiece)
    move.oldPiece.set(toPiece)

    move.location.clear()
    move.location.set(fromLocation)
    move.location.set(toLocation)

    move.newTarget = nil
    move.newRisk = nil

    move.oldTarget = target
    move.oldRisk = risk

    move.newStart = nil
    move.newEnd = nil

    move.oldStart = start
    move.oldEnd = end

    moves.next()
}

func addMoveRevealEnPassant(
    b *SimpleBoard,
    fromPiece *Piece,
    fromLocation *Point,
    toPiece *Piece,
    toLocation *Point,
    newPiece *Piece,
    newTarget *Point,
    newRisk *Point,
    moves *Array1000[FastMove],
) {
    move := moves.get()

    addMoveSimple(b, fromPiece, fromLocation, toPiece, toLocation, newPiece, moves)

    move.newTarget = newTarget
    move.newRisk = newRisk
}

func addMoveCaptureEnPassant(
    b *SimpleBoard,
    fromPiece *Piece,
    fromLocation *Point,
    toPiece *Piece,
    toLocation *Point,
    newPiece *Piece,
    risk1 *Point,
    risk2 *Point,
    moves *Array1000[FastMove],
) {
    move := moves.get()

    addMoveSimple(b, fromPiece, fromLocation, toPiece, toLocation, newPiece, moves)

    capturedPiece := b.getPiece(risk1)
    if capturedPiece == nil {
        return
    }

    move.newPiece.set(nil)

    move.oldPiece.set(capturedPiece)

    move.location.set(risk1)

    capturedPiece = b.getPiece(risk2)
    if capturedPiece == nil {
        return
    }

    move.newPiece.set(nil)

    move.oldPiece.set(capturedPiece)

    move.location.set(risk2)
}

func addMoveAllyDefense(
    b *SimpleBoard,
    fromPiece *Piece,
    fromLocation *Point,
    toLocation *Point,
    moves *Array1000[FastMove],
) {
    move := moves.get()

    move.b = b
    move.fromLocation = fromLocation
    move.toLocation = toLocation
    move.color = fromPiece.color
    move.allyDefense = true
    move.promotionIndex = -1

    move.newPiece.clear()
    move.oldPiece.clear()
    move.location.clear()

    move.newTarget = nil
    move.newRisk = nil
    move.oldTarget = nil
    move.oldRisk = nil

    move.newStart = nil
    move.newEnd = nil
    move.oldStart = nil
    move.oldEnd = nil

    moves.next()
}

func addMoveCastle(
    b *SimpleBoard,
    king *Piece,
    fromLocation *Point,
    toKingLocation *Point,
    rook *Piece,
    toLocation *Point,
    toRookLocation *Point,
    newStart *Point,
    newEnd *Point,
    moves *Array1000[FastMove],
) {
    move := moves.get()

    move.b = b
    move.fromLocation = fromLocation
    move.toLocation = toLocation
    move.color = king.color
    move.allyDefense = false
    move.promotionIndex = -1

    target, risk := b.getEnPassant(king.color)
    start, end := b.getVulnerable(king.color)

    newKing := b.movePiece(king)
    newRook := b.movePiece(rook)

    move.newPiece.clear()
    move.newPiece.set(nil)
    move.newPiece.set(nil)
    move.newPiece.set(newKing)
    move.newPiece.set(newRook)

    move.oldPiece.clear()
    move.oldPiece.set(king)
    move.oldPiece.set(rook)
    move.oldPiece.set(nil)
    move.oldPiece.set(nil)

    move.location.clear()
    move.location.set(fromLocation)
    move.location.set(toLocation)
    move.location.set(toKingLocation)
    move.location.set(toRookLocation)

    move.newTarget = nil
    move.newRisk = nil

    move.oldTarget = target
    move.oldRisk = risk

    move.newStart = newStart
    move.newEnd = newEnd

    move.oldStart = start
    move.oldEnd = end

    moves.next()
}

type FastMove struct {
    b *SimpleBoard
    fromLocation *Point
    toLocation *Point
    color int
    allyDefense bool
    promotionIndex int

    // piece changes
    newPiece Array4[*Piece]
    oldPiece Array4[*Piece]
    location Array4[*Point]

    // enPassant
    newTarget *Point
    oldTarget *Point
    newRisk *Point
    oldRisk *Point

    // vulnerable
    newStart *Point
    newEnd *Point
    oldStart *Point
    oldEnd *Point
}

func (m *FastMove) execute() {
    for i := 0; i < m.newPiece.count; i++ {
        m.b.setPiece(m.location.array[i], m.newPiece.array[i])
    }

    m.b.setEnPassant(m.color, m.newTarget, m.newRisk)
    m.b.setVulnerable(m.color, m.newStart, m.newEnd)
}

func (m *FastMove) undo() {
    for i := m.newPiece.count - 1; i >= 0; i-- {
        m.b.setPiece(m.location.array[i], m.oldPiece.array[i])
    }

    m.b.setEnPassant(m.color, m.oldTarget, m.oldRisk)
    m.b.setVulnerable(m.color, m.oldStart, m.oldEnd)
}

