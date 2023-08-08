class Move {
    constructor() {}

    execute() {
        throw new Error('Move.execute not implemented')
    }

    undo() {
        throw new Error('Move.undo not implemented')
    }
}

class SimpleMove extends Move {
    constructor(board, fromX, fromY, toX, toY, options) {
        super()
        this.board = board

        this.fromX = fromX
        this.fromY = fromY
        this.toX = toX
        this.toY = toY

        this.piece = board.getPiece(fromX, fromY)
        this.newPiece = this.piece.copy()
        this.newPiece.move()
        this.capturedPiece = board.getPiece(toX, toY)
    }

    execute() {
        this.board.setPiece(this.toX, this.toY, this.newPiece)
        this.board.setPiece(this.fromX, this.fromY, null)
        this.board.increment()

        return true
    }

    undo() {
        this.board.setPiece(this.fromX, this.fromY, this.piece)
        this.board.setPiece(this.toX, this.toY, this.capturedPiece)
        this.board.decrement()

        return true
    }
}

// class CastleMove extends Move {
//     constructor(board, kFromX, kFromY, kToX, kToY, rFromX, rFromY, rToX, rToY) {
//         super()
//         this.board = board
//         this.fromX = kFromX
//         this.fromY = kFromY
//         this.toX = kToX
//         this.toY = kToY
//         this.rookFromX = rFromX
//         this.rookFromY = rFromY
//         this.rookToX = rToX
//         this.rookToY = rToY
//         this.king = board.getPiece(kFromX, kFromY)
//         this.rook = board.getPiece(rFromX, rFromY)
//         this.kingCapturedPiece = board.getPiece(kToX, kToY)
//         this.rookCapturedPiece = board.getPiece(rToX, rToY)
//     }
// }

// class PromotionMove extends Move {
//     constructor(board, fromX, fromY, toX, toY, promotionPiece) {
//         super()
//         this.board = board
//         this.fromX = fromX
//         this.fromY = fromY
//         this.toX = toX
//         this.toY = toY
//         this.promotionPiece = promotionPiece
//         this.piece = board.getPiece(fromX, fromY)
//         this.capturedPiece = board.getPiece(toX, toY)
//     }
// }

// class EnPassantMove extends Move {
//     constructor(board, fromX, fromY, toX, toY) {
//         super()
//         this.board = board
//         this.fromX = fromX
//         this.fromY = fromY
//         this.toX = toX
//         this.toY = toY
//         this.piece = board.getPiece(fromX, fromY)
//         this.capturedPiece = board.getPiece(toX, toY)
//     }
// }

export { Move, SimpleMove }
