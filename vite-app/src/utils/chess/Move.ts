import { Piece } from './Piece'
import { Board } from './Board'
import { BoardMove, EnPassant } from './Types'

class MoveFactory {
    static createMove(board: Board, m: BoardMove): Move | null {
        try {
            if (m.options.revealEnPassant) {
                return new RevealEnPassantMove(board, m)
            } else if (m.options.captureEnPassant) {
                return new CaptureEnPassantMove(board, m)
            } else {
                return new SimpleMove(board, m)
            }
            // TODO: castle move
            // TODO: promote move
        } catch (e) {
            console.error(e)
            return null
        }
    }
}

abstract class Move {
    abstract execute(): boolean
    abstract undo(): boolean
}

class SimpleMove extends Move {
    protected board: Board
    protected m: BoardMove
    protected piece: Piece
    protected newPiece: Piece
    protected capturedPiece: Piece | null
    protected enPassant: EnPassant | null

    constructor(board: Board, m: BoardMove) {
        super()
        this.board = board
        this.m = m

        const tempPiece = board.getPiece(m.xFrom, m.yFrom)
        if (!tempPiece) {
            throw new Error('Invalid from position')
        }

        const tempCapturedPiece = board.getPiece(m.xTo, m.yTo)
        if (tempCapturedPiece === undefined) {
            throw new Error('Invalid to position')
        }

        this.piece = tempPiece
        this.newPiece = this.piece.copy()
        this.capturedPiece = tempCapturedPiece
        this.enPassant = board.getEnPassant(this.piece.color)
    }

    execute(): boolean {
        this.board.setPiece(this.m.xTo, this.m.yTo, this.newPiece)
        this.board.setPiece(this.m.xFrom, this.m.yFrom, null)
        this.board.setEnPassant(this.piece.color, null)
        this.board.increment()

        return true
    }

    undo(): boolean {
        this.board.setPiece(this.m.xFrom, this.m.yFrom, this.piece)
        this.board.setPiece(this.m.xTo, this.m.yTo, this.capturedPiece)
        this.board.setEnPassant(this.piece.color, this.enPassant)
        this.board.decrement()

        return true
    }
}

class RevealEnPassantMove extends SimpleMove {
    private revealEnPassant: {x: number, y: number}

    constructor(board: Board, m: BoardMove) {
        super(board, m)

        if (!m.options.revealEnPassant) {
            throw new Error('Invalid move options')
        }

        this.revealEnPassant = m.options.revealEnPassant
    }

    execute(): boolean {
        const res = super.execute()
        this.board.setEnPassant(
            this.piece.color,
            {
                x: this.revealEnPassant.x,
                y: this.revealEnPassant.y,
                xPiece: this.m.xTo,
                yPiece: this.m.yTo,
            }
        )
        return res
    }

    undo(): boolean {
        return super.undo()
    }
}

class CaptureEnPassantMove extends SimpleMove {
    private enPassantsAndCapturedPieces: {enPassant: EnPassant, capturedPiece: Piece}[]

    constructor(board: Board, m: BoardMove) {
        super(board, m)

        if (!m.options.captureEnPassant) {
            throw new Error('Invalid move options')
        }

        this.enPassantsAndCapturedPieces = this.board.getEnPassants(this.piece.color, m.xTo, m.yTo)
        .map((e) => ({enPassant: e, capturedPiece: board.getPiece(e.xPiece, e.yPiece)?.copy()}))
        .filter((e): e is {enPassant: EnPassant, capturedPiece: Piece} => true)
    }

    execute() {
        const res = super.execute()
        this.enPassantsAndCapturedPieces.forEach((e) => {
            this.board.setPiece(e.enPassant.xPiece, e.enPassant.yPiece, null)
        })
        return res
    }

    undo() {
        const res = super.undo()
        this.enPassantsAndCapturedPieces.forEach((e) => {
            this.board.setPiece(e.enPassant.xPiece, e.enPassant.yPiece, e.capturedPiece)
        })
        return res
    }
}

export { MoveFactory, Move, SimpleMove, RevealEnPassantMove, CaptureEnPassantMove }