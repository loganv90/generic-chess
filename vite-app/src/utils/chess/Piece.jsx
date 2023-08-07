class Piece {
    constructor(color) {
        this.color = color
    }

    getColor() { return this.color }

    move() {}

    getEnPassant() { return {} }

    setEnPassant() {}

    copy() {
        throw new Error('Piece.copy not implemented')
    }

    getMoves() {
        throw new Error('Piece.getMoves not implemented')
    }
}

class Pawn extends Piece {
    constructor(color, xDir, yDir, moved=false, enPassant={}) {
        super(color)
        this.moved = moved
        this.enPassant = enPassant

        this.xDir = xDir > 0 ? 1 : xDir < 0 ? -1 : 0
        this.yDir = yDir > 0 ? 1 : yDir < 0 ? -1 : 0

        if (this.xDir && this.yDir) {
            throw new Error('Pawn can not move diagonally')
        } else if (!this.xDir && !this.yDir) {
            throw new Error('Pawn must have a direction')
        }
    }

    getEnPassant() {
        return this.enPassant
    }

    setEnPassant(enPassant) {
        this.enPassant = enPassant
    }

    copy() {
        return new Pawn(this.color, this.xDir, this.yDir, this.moved, this.enPassant)
    }

    move() {
        this.moved = true
    }

    getMoves() {
        const moves = [{
            x: this.xDir,
            y: this.yDir,
            options: {
                noCapture: true,
                canPromote: true
            }
        }]

        if (!this.moved) {
            moves.push({
                x: this.xDir*2,
                y: this.yDir*2,
                options: {
                    noCapture: true,
                    canPromote: true,
                    enPassant: {
                        x: this.xDir,
                        y: this.yDir
                    }
                }
            })
        }

        if (this.xDir) {
            moves.push({
                x: this.xDir,
                y: 1,
                options: {
                    mustCapture: true,
                    canPromote: true,
                    canEnPassant: true
                }
            })
            moves.push({
                x: this.xDir,
                y: -1,
                options: {
                    mustCapture: true,
                    canPromote: true,
                    canEnPassant: true
                }
            })
        } else if (this.yDir) {
            moves.push({
                x: 1,
                y: this.yDir,
                options: {
                    mustCapture: true,
                    canPromote: true,
                    canEnPassant: true
                }
            })
            moves.push({
                x: -1,
                y: this.yDir,
                options: {
                    mustCapture: true,
                    canPromote: true,
                    canEnPassant: true
                }
            })
        } else {
            throw new Error('Pawn does not have a direction')
        }

        return moves
    }
}

class Knight extends Piece {
    constructor(color) {
        super(color)
        this.options = {}
    }

    copy() {
        return new Knight(this.color)
    }

    getMoves() {
        return [
            {x: 1, y: 2, options: this.options},
            {x: 1, y: -2, options: this.options},
            {x: -1, y: -2, options: this.options},
            {x: 2, y: 1, options: this.options},
            {x: 2, y: -1, options: this.options},
            {x: -2, y: 1, options: this.options},
            {x: -2, y: -1, options: this.options},
        ]
    }
}

class Bishop extends Piece {
    constructor(color) {
        super(color)
        this.options = {isDirection: true}
    }

    copy() {
        return new Bishop(this.color)
    }

    static getMoves() {
        return [
            {x: 1, y: 1, options: this.options},
            {x: 1, y: -1, options: this.options},
            {x: -1, y: 1, options: this.options},
            {x: -1, y: -1, options: this.options},
        ]
    }
}

class Rook extends Piece {
    constructor(color, moved=false) {
        super(color)
        this.moved = moved
        this.options = {isDirection: true}
    }

    copy() {
        return new Rook(this.color, this.moved)
    }

    move() {
        this.moved = true
    }

    static getMoves() {
        return [
            {x: 1, y: 0, options: this.options},
            {x: -1, y: 0, options: this.options},
            {x: 0, y: 1, options: this.options},
            {x: 0, y: -1, options: this.options},
        ]
    }
}

class Queen extends Piece {
    constructor(color) {
        super(color)
        this.options = {isDirection: true}
    }

    copy() {
        return new Queen(this.color)
    }

    static getMoves() {
        return [
            {x: 1, y: 1, options: this.options},
            {x: 1, y: -1, options: this.options},
            {x: -1, y: 1, options: this.options},
            {x: -1, y: -1, options: this.options},
            {x: 1, y: 0, options: this.options},
            {x: -1, y: 0, options: this.options},
            {x: 0, y: 1, options: this.options},
            {x: 0, y: -1, options: this.options},
        ]
    }
}

class King extends Piece {
    constructor(color, moved=false) {
        super(color)
        this.moved = moved
        this.options = {}
    }

    copy() {
        return new King(this.color, this.moved)
    }

    move() {
        this.moved = true
    }

    static getMoves() {
        return [
            {x: 1, y: 1, options: this.options},
            {x: 1, y: -1, options: this.options},
            {x: -1, y: 1, options: this.options},
            {x: -1, y: -1, options: this.options},
            {x: 1, y: 0, options: this.options},
            {x: -1, y: 0, options: this.options},
            {x: 0, y: 1, options: this.options},
            {x: 0, y: -1, options: this.options},
        ]
    }
}

export { Piece, Pawn, Knight, Bishop, Rook, Queen, King }
