package chess

func addMoveSimple(
    b *SimpleBoard,
    fromPiece *Piece,
    fromLocation *Point,
    toPiece *Piece,
    toLocation *Point,
    newPiece *Piece,
) {
    var move *FastMove
    color := fromPiece.color

    if toPiece != nil { // capture
        move = b.captureMoves[color].get()
        move.captureValue = toPiece.value()
    } else { // no capture
        move = b.moves[color].get()
        move.captureValue = 0
    }

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
) {
    var move *FastMove
    color := fromPiece.color

    if toPiece != nil { // capture
        move = b.captureMoves[color].get()
        move.captureValue = toPiece.value()
    } else { // no capture
        move = b.moves[color].get()
        move.captureValue = 0
    }

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

    move.newTarget = newTarget
    move.newRisk = newRisk

    move.oldTarget = target
    move.oldRisk = risk

    move.newStart = nil
    move.newEnd = nil

    move.oldStart = start
    move.oldEnd = end
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
) {
    var move *FastMove
    color := fromPiece.color
    capturedPiece1 := b.getPiece(risk1)
    capturedPiece2 := b.getPiece(risk2)

    if toPiece != nil || capturedPiece1 == nil || capturedPiece2 == nil { // capture
        move = b.captureMoves[color].get()
    } else {
        move = b.moves[color].get()
    }

    move.b = b
    move.fromLocation = fromLocation
    move.toLocation = toLocation
    move.color = fromPiece.color
    move.allyDefense = false
    move.captureValue = 0

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

    if toPiece != nil {
        move.captureValue += toPiece.value()
    }
    if capturedPiece1 != nil {
        move.newPiece.set(nil)
        move.oldPiece.set(capturedPiece1)
        move.location.set(risk1)
        move.captureValue += capturedPiece1.value()
    }
    if capturedPiece2 != nil {
        move.newPiece.set(nil)
        move.oldPiece.set(capturedPiece2)
        move.location.set(risk2)
        move.captureValue += capturedPiece2.value()
    }
}

func addMoveAllyDefense(
    b *SimpleBoard,
    fromPiece *Piece,
    fromLocation *Point,
    toLocation *Point,
) {
    var move *FastMove
    color := fromPiece.color

    move = b.defenseMoves[color].get()

    move.b = b
    move.fromLocation = fromLocation
    move.toLocation = toLocation
    move.color = fromPiece.color
    move.allyDefense = true
    move.promotionIndex = -1
    move.captureValue = 0
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
) {
    var move *FastMove
    color := king.color

    move = b.moves[color].get()

    move.b = b
    move.fromLocation = fromLocation
    move.toLocation = toLocation
    move.color = king.color
    move.allyDefense = false
    move.promotionIndex = -1
    move.captureValue = 0

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
}

type FastMove struct {
    b *SimpleBoard
    fromLocation *Point
    toLocation *Point
    color int
    allyDefense bool
    promotionIndex int
    captureValue int

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

