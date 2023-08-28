class MoveFactory {
    static createMove(board, fromX, fromY, toX, toY, options) {
        if (options?.enPassant) {
            return new RevealEnPassantMove(
                board,
                fromX,
                fromY,
                toX,
                toY,
                options.enPassant.x,
                options.enPassant.y
            )
        } else if (options?.canEnPassant) {
            return new CaptureEnPassantMove(
                board,
                fromX,
                fromY,
                toX,
                toY
            )
        } else {
            return new SimpleMove(
                board,
                fromX,
                fromY,
                toX,
                toY
            )
        }
    }
}

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
    constructor(board, fromX, fromY, toX, toY) {
        super()
        this.board = board

        this.fromX = fromX
        this.fromY = fromY
        this.toX = toX
        this.toY = toY

        this.piece = board.getPiece(fromX, fromY)
        this.enPassant = board.getEnPassant(this.piece.color)
        this.newPiece = this.piece.copy()
        this.newPiece.move()
        this.capturedPiece = board.getPiece(toX, toY)
    }

    execute() {
        this.board.setPiece(this.toX, this.toY, this.newPiece)
        this.board.setPiece(this.fromX, this.fromY, null)
        this.board.increment()
        this.board.setEnPassant(this.piece.color, null)

        return true
    }

    undo() {
        this.board.setPiece(this.fromX, this.fromY, this.piece)
        this.board.setPiece(this.toX, this.toY, this.capturedPiece)
        this.board.decrement()
        this.board.setEnPassant(this.piece.color, this.enPassant)

        return true
    }
}

class RevealEnPassantMove extends SimpleMove {
    constructor(board, fromX, fromY, toX, toY, enPassantX, enPassantY) {
        super(board, fromX, fromY, toX, toY)

        this.enPassantX = enPassantX
        this.enPassantY = enPassantY
    }

    execute() {
        const res = super.execute()
        this.board.setEnPassant(
            this.piece.color,
            {
                x: this.enPassantX,
                y: this.enPassantY,
                xPiece: this.toX,
                yPiece: this.toY
            }
        )
        return res
    }

    undo() {
        return super.undo()
    }
}

class CaptureEnPassantMove extends SimpleMove {
    constructor(board, fromX, fromY, toX, toY) {
        super(board, fromX, fromY, toX, toY)

        this.enPassantsAndPieces = this.board.getEnPassants(this.piece.color, toX, toY)
        .map((e) => { return {enPassant: e, piece: board.getPiece(e.xPiece, e.yPiece).copy()} })
    }

    execute() {
        const res = super.execute()
        this.enPassantsAndPieces.forEach((e) => {
            this.board.setPiece(e.enPassant.xPiece, e.enPassant.yPiece, null)
        })
        return res
    }

    undo() {
        const res = super.undo()
        this.enPassantsAndPieces.forEach((e) => {
            this.board.setPiece(e.enPassant.xPiece, e.enPassant.yPiece, e.piece)
        })
        return res
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

export { MoveFactory, Move, SimpleMove, RevealEnPassantMove, CaptureEnPassantMove }
