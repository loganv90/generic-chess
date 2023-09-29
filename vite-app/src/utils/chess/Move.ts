import { Piece } from './Piece'
import { Board } from './Board'

abstract class Move {
    abstract execute(): boolean
    abstract undo(): boolean
}

class SimpleMove extends Move {
    protected readonly board: Board
    protected readonly xFrom: number
    protected readonly yFrom: number
    protected readonly xTo: number
    protected readonly yTo: number
    protected readonly piece: Piece
    protected readonly newPiece: Piece
    protected readonly capturedPiece: Piece | null

    constructor(
        board: Board,
        xFrom: number,
        yFrom: number,
        xTo: number,
        yTo: number,
    ) {
        super()
        this.board = board
        this.xFrom = xFrom
        this.yFrom = yFrom
        this.xTo = xTo
        this.yTo = yTo

        const tempPiece = board.getPiece(xFrom, yFrom)
        if (!tempPiece) {
            throw new Error('Invalid from position')
        }

        const tempCapturedPiece = board.getPiece(xTo, yTo)
        if (tempCapturedPiece === undefined) {
            throw new Error('Invalid to position')
        }

        this.piece = tempPiece
        this.newPiece = this.piece.copy()
        this.capturedPiece = tempCapturedPiece
    }

    execute(): boolean {
        this.board.setPiece(this.xTo, this.yTo, this.newPiece)
        this.board.setPiece(this.xFrom, this.yFrom, null)
        this.board.increment()

        return true
    }

    undo(): boolean {
        this.board.setPiece(this.xFrom, this.yFrom, this.piece)
        this.board.setPiece(this.xTo, this.yTo, this.capturedPiece)
        this.board.decrement()

        return true
    }
}

// class RevealEnPassantMove extends SimpleMove {
//     private revealEnPassant: {x: number, y: number}

//     constructor(board: Board, m: BoardMove) {
//         super(board, m)

//         if (!m.options.revealEnPassant) {
//             throw new Error('Invalid move options')
//         }

//         this.revealEnPassant = m.options.revealEnPassant
//     }

//     execute(): boolean {
//         const res = super.execute()
//         this.board.setEnPassant(
//             this.piece.color,
//             {
//                 x: this.revealEnPassant.x,
//                 y: this.revealEnPassant.y,
//                 xPiece: this.m.xTo,
//                 yPiece: this.m.yTo,
//             }
//         )
//         return res
//     }

//     undo(): boolean {
//         return super.undo()
//     }
// }

// class CaptureEnPassantMove extends SimpleMove {
//     private enPassantsAndCapturedPieces: {enPassant: EnPassant, capturedPiece: Piece}[]

//     constructor(board: Board, m: BoardMove) {
//         super(board, m)

//         if (!m.options.captureEnPassant) {
//             throw new Error('Invalid move options')
//         }

//         this.enPassantsAndCapturedPieces = this.board.getEnPassants(this.piece.color, m.xTo, m.yTo)
//         .map((e) => ({enPassant: e, capturedPiece: board.getPiece(e.xPiece, e.yPiece)?.copy()}))
//         .filter((e): e is {enPassant: EnPassant, capturedPiece: Piece} => true)
//     }

//     execute() {
//         const res = super.execute()
//         this.enPassantsAndCapturedPieces.forEach((e) => {
//             this.board.setPiece(e.enPassant.xPiece, e.enPassant.yPiece, null)
//         })
//         return res
//     }

//     undo() {
//         const res = super.undo()
//         this.enPassantsAndCapturedPieces.forEach((e) => {
//             this.board.setPiece(e.enPassant.xPiece, e.enPassant.yPiece, e.capturedPiece)
//         })
//         return res
//     }
// }

export { Move, SimpleMove } 