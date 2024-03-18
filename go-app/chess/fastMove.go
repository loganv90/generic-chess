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
    move.newPiece.next()
    move.newPiece.set(newPiece)
    move.newPiece.next()

    move.newLocation.clear()
    move.newLocation.set(fromLocation)
    move.newLocation.next()
    move.newLocation.set(toLocation)
    move.newLocation.next()

    move.oldPiece.clear()
    move.oldPiece.set(fromPiece)
    move.oldPiece.next()
    move.oldPiece.set(toPiece)
    move.oldPiece.next()

    move.oldLocation.clear()
    move.oldLocation.set(fromLocation)
    move.oldLocation.next()
    move.oldLocation.set(toLocation)
    move.oldLocation.next()

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
    target1 *Point,
    target2 *Point,
    risk1 *Point,
    risk2 *Point,
    moves *Array1000[FastMove],
) {
    move := moves.get()

    addMoveSimple(b, fromPiece, fromLocation, toPiece, toLocation, newPiece, moves)

    if target1 == nil {
        return
    }

    capturedPiece := b.getPiece(risk1)
    if capturedPiece == nil {
        return
    }

    move.newPiece.set(nil)
    move.newPiece.next()

    move.newLocation.set(risk1)
    move.newLocation.next()

    move.oldPiece.set(capturedPiece)
    move.oldPiece.next()

    move.oldLocation.set(risk1)
    move.oldLocation.next()

    if target2 == nil {
        return
    }

    capturedPiece = b.getPiece(risk2)
    if capturedPiece == nil {
        return
    }

    move.newPiece.set(nil)
    move.newPiece.next()

    move.newLocation.set(risk2)
    move.newLocation.next()

    move.oldPiece.set(capturedPiece)
    move.oldPiece.next()

    move.oldLocation.set(risk2)
    move.oldLocation.next()
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
    move.newLocation.clear()
    move.oldPiece.clear()
    move.oldLocation.clear()

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
    move.newPiece.next()
    move.newPiece.set(nil)
    move.newPiece.next()
    move.newPiece.set(newKing)
    move.newPiece.next()
    move.newPiece.set(newRook)
    move.newPiece.next()

    move.newLocation.clear()
    move.newLocation.set(fromLocation)
    move.newLocation.next()
    move.newLocation.set(toLocation)
    move.newLocation.next()
    move.newLocation.set(toKingLocation)
    move.newLocation.next()
    move.newLocation.set(toRookLocation)
    move.newLocation.next()

    move.oldPiece.clear()
    move.oldPiece.set(nil)
    move.oldPiece.next()
    move.oldPiece.set(nil)
    move.oldPiece.next()
    move.oldPiece.set(king)
    move.oldPiece.next()
    move.oldPiece.set(rook)
    move.oldPiece.next()

    move.oldLocation.clear()
    move.oldLocation.set(toKingLocation)
    move.oldLocation.next()
    move.oldLocation.set(toRookLocation)
    move.oldLocation.next()
    move.oldLocation.set(fromLocation)
    move.oldLocation.next()
    move.oldLocation.set(toLocation)
    move.oldLocation.next()

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
    newLocation Array4[*Point]
    oldPiece Array4[*Piece]
    oldLocation Array4[*Point]

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
        m.b.setPiece(m.newLocation.array[i], m.newPiece.array[i])
    }

    m.b.setEnPassant(m.color, m.newTarget, m.newRisk)
    m.b.setVulnerable(m.color, m.newStart, m.newEnd)
}

func (m *FastMove) undo() {
    for i := 0; i < m.oldPiece.count; i++ {
        m.b.setPiece(m.oldLocation.array[i], m.oldPiece.array[i])
    }

    m.b.setEnPassant(m.color, m.oldTarget, m.oldRisk)
    m.b.setVulnerable(m.color, m.oldStart, m.oldEnd)
}

