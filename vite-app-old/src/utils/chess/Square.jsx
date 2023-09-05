import { Pawn, Knight, Bishop, Rook, Queen, King } from '/src/utils/chess/Piece'

class Square {
    constructor(id, x, y, isLight, fenChar='') {
        this.id = id
        this.x = x
        this.y = y
        this.isLight = isLight
        this.piece = this.createPieceFromFen(fenChar)
    }

    getMoves() {
        return this.piece ? this.piece.getMoves() : []
    }

    createPieceFromFen(fenChar) {
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
