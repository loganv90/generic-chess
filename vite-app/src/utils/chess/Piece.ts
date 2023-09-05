import { PieceMove } from './Types'

abstract class Piece {
    readonly color: string
    readonly moved: boolean

    constructor(color: string, moved=false) {
        this.color = color
        this.moved = moved
    }

    abstract copy(): Piece
    abstract getMoves(): PieceMove[]
}

class Pawn extends Piece {
    private xDir: 1|0|-1
    private yDir: 1|0|-1

    constructor(color: string, xDir: 1|0|-1, yDir: 1|0|-1, moved=false) {
        super(color, moved)

        this.xDir = xDir
        this.yDir = yDir

        if (this.xDir && this.yDir) {
            throw new Error('Pawn can not move diagonally')
        } else if (!this.xDir && !this.yDir) {
            throw new Error('Pawn must have a direction')
        }
    }

    copy(): Pawn {
        return new Pawn(this.color, this.xDir, this.yDir, true)
    }

    getMoves(): PieceMove[] {
        const moves: PieceMove[] = [{x: this.xDir, y: this.yDir, options: {
            noCapture: true,
            canPromote: true,
        }}]

        if (!this.moved) {
            moves.push({x: this.xDir*2, y: this.yDir*2, options: {
                noCapture: true,
                canPromote: true,
                mustRevealEnPassant: {x: this.xDir, y: this.yDir},
            }})
        }

        if (this.xDir) {
            moves.push({x: this.xDir, y: 1, options: {
                mustCapture: true,
                canPromote: true,
                canCaptureEnPassant: true,
            }})
            moves.push({x: this.xDir, y: -1, options: {
                mustCapture: true,
                canPromote: true,
                canCaptureEnPassant: true,
            }})
        } else if (this.yDir) {
            moves.push({x: 1, y: this.yDir, options: {
                mustCapture: true,
                canPromote: true,
                canCaptureEnPassant: true,
            }})
            moves.push({x: -1, y: this.yDir, options: {
                mustCapture: true,
                canPromote: true,
                canCaptureEnPassant: true,
            }})
        } else {
            throw new Error('Pawn does not have a direction')
        }

        return moves
    }
}

class Knight extends Piece {
    constructor(color: string) {
        super(color)
    }

    copy(): Knight {
        return new Knight(this.color)
    }

    getMoves(): PieceMove[] {
        return [
            {x: 1, y: 2, options: {}},
            {x: 1, y: -2, options: {}},
            {x: -1, y: -2, options: {}},
            {x: -1, y: 2, options: {}},
            {x: 2, y: 1, options: {}},
            {x: 2, y: -1, options: {}},
            {x: -2, y: 1, options: {}},
            {x: -2, y: -1, options: {}},
        ]
    }
}

class Bishop extends Piece {
    constructor(color: string) {
        super(color)
    }

    copy(): Bishop {
        return new Bishop(this.color)
    }

    getMoves(): PieceMove[] {
        return [
            {x: 1, y: 1, options: {direction: true}},
            {x: 1, y: -1, options: {direction: true}},
            {x: -1, y: 1, options: {direction: true}},
            {x: -1, y: -1, options: {direction: true}},
        ]
    }
}

class Rook extends Piece {
    constructor(color: string, moved=false) {
        super(color, moved)
    }

    copy(): Rook {
        return new Rook(this.color, true)
    }

    getMoves(): PieceMove[] {
        return [
            {x: 1, y: 0, options: {direction: true}},
            {x: -1, y: 0, options: {direction: true}},
            {x: 0, y: 1, options: {direction: true}},
            {x: 0, y: -1, options: {direction: true}},
        ]
    }
}

class Queen extends Piece {
    constructor(color: string) {
        super(color)
    }

    copy(): Queen {
        return new Queen(this.color)
    }

    getMoves(): PieceMove[] {
        return [
            {x: 1, y: 1, options: {direction: true}},
            {x: 1, y: -1, options: {direction: true}},
            {x: -1, y: 1, options: {direction: true}},
            {x: -1, y: -1, options: {direction: true}},
            {x: 1, y: 0, options: {direction: true}},
            {x: -1, y: 0, options: {direction: true}},
            {x: 0, y: 1, options: {direction: true}},
            {x: 0, y: -1, options: {direction: true}},
        ]
    }
}

class King extends Piece {
    constructor(color: string, moved=false) {
        super(color, moved)
    }

    copy(): King {
        return new King(this.color, true)
    }

    getMoves(): PieceMove[] {
        return [
            {x: 1, y: 1, options: {}},
            {x: 1, y: -1, options: {}},
            {x: -1, y: 1, options: {}},
            {x: -1, y: -1, options: {}},
            {x: 1, y: 0, options: {}},
            {x: -1, y: 0, options: {}},
            {x: 0, y: 1, options: {}},
            {x: 0, y: -1, options: {}},
        ]
    }
}

export { Piece, Pawn, Knight, Bishop, Rook, Queen, King }