class Piece {
    constructor(color) {
        this.color = color
    }

    move() {}

    copy() {
        throw new Error('Piece.copy not implemented')
    }

    getMoves() {
        throw new Error('Piece.getMoves not implemented')
    }
}

class Pawn extends Piece {
    constructor(color, xDir, yDir, moved=false) {
        super(color)
        this.moved = moved

        this.xDir = xDir > 0 ? 1 : xDir < 0 ? -1 : 0
        this.yDir = yDir > 0 ? 1 : yDir < 0 ? -1 : 0

        if (this.xDir && this.yDir) {
            throw new Error('Pawn can not move diagonally')
        } else if (!this.xDir && !this.yDir) {
            throw new Error('Pawn must have a direction')
        }
    }

    copy() {
        return new Pawn(this.color, this.xDir, this.yDir, this.moved)
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
                    mustCross: {
                        x: this.xDir,
                        y: this.yDir
                    },
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
    }

    copy() {
        return new Knight(this.color)
    }

    getMoves() {
        return [
            {x: 1, y: 2},
            {x: 1, y: -2},
            {x: -1, y: -2},
            {x: -1, y: 2},
            {x: 2, y: 1},
            {x: 2, y: -1},
            {x: -2, y: 1},
            {x: -2, y: -1},
        ]
    }
}

class Bishop extends Piece {
    constructor(color) {
        super(color)
    }

    copy() {
        return new Bishop(this.color)
    }

    getMoves() {
        return [
            {x: 1, y: 1, options: {isDirection: true}},
            {x: 1, y: -1, options: {isDirection: true}},
            {x: -1, y: 1, options: {isDirection: true}},
            {x: -1, y: -1, options: {isDirection: true}},
        ]
    }
}

class Rook extends Piece {
    constructor(color, moved=false) {
        super(color)
        this.moved = moved
    }

    copy() {
        return new Rook(this.color, this.moved)
    }

    move() {
        this.moved = true
    }

    getMoves() {
        return [
            {x: 1, y: 0, options: {isDirection: true}},
            {x: -1, y: 0, options: {isDirection: true}},
            {x: 0, y: 1, options: {isDirection: true}},
            {x: 0, y: -1, options: {isDirection: true}},
        ]
    }
}

class Queen extends Piece {
    constructor(color) {
        super(color)
    }

    copy() {
        return new Queen(this.color)
    }

    getMoves() {
        return [
            {x: 1, y: 1, options: {isDirection: true}},
            {x: 1, y: -1, options: {isDirection: true}},
            {x: -1, y: 1, options: {isDirection: true}},
            {x: -1, y: -1, options: {isDirection: true}},
            {x: 1, y: 0, options: {isDirection: true}},
            {x: -1, y: 0, options: {isDirection: true}},
            {x: 0, y: 1, options: {isDirection: true}},
            {x: 0, y: -1, options: {isDirection: true}}
        ]
    }
}

class King extends Piece {
    constructor(color, moved=false) {
        super(color)
        this.moved = moved
    }

    copy() {
        return new King(this.color, this.moved)
    }

    move() {
        this.moved = true
    }

    getMoves() {
        return [
            {x: 1, y: 1},
            {x: 1, y: -1},
            {x: -1, y: 1},
            {x: -1, y: -1},
            {x: 1, y: 0},
            {x: -1, y: 0},
            {x: 0, y: 1},
            {x: 0, y: -1}
        ]
    }
}

export { Piece, Pawn, Knight, Bishop, Rook, Queen, King }
