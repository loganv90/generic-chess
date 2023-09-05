import { PieceMove } from './Types'
import { Piece, Pawn, Knight, Bishop, Rook, Queen, King } from './Piece'

class Square {
    readonly id: string
    readonly x: number
    readonly y: number
    readonly isLight: boolean
    private piece: Piece | null

    constructor(id: string, x: number, y: number, fenChar='') {
        this.id = id
        this.x = x
        this.y = y
        this.isLight = (x+y) % 2 == 0
        this.piece = this.createPieceFromFen(fenChar)
    }

    getPiece(): Piece | null {
        return this.piece
    }

    setPiece(piece: Piece | null): void {
        this.piece = piece
    }

    getMoves(): PieceMove[] {
        return this.piece ? this.piece.getMoves() : []
    }

    private createPieceFromFen(fenChar: string): Piece | null {
        switch (fenChar) {
            case 'r':
                return new Rook('black')
            case 'n':
                return new Knight('black')
            case 'b':
                return new Bishop('black')
            case 'q':
                return new Queen('black')
            case 'k':
                return new King('black')
            case 'p':
                return new Pawn('black', 0, 1)
            case 'R':
                return new Rook('white')
            case 'N':
                return new Knight('white')
            case 'B':
                return new Bishop('white')
            case 'Q':
                return new Queen('white')
            case 'K':
                return new King('white')
            case 'P':
                return new Pawn('white', 0, -1)
            default:
                return null
        }
    }
}

export { Square }