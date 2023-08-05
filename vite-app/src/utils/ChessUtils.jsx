function createPieceFromFen(c) {
    switch (c) {
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

function createSquaresFromFen(fenPieces) {
    const squareRows = []
    const rows = fenPieces.split('/')

    for (const row of rows) {
        const squares = []

        for (const c of row) {
            let x = squares.length
            let y = squareRows.length

            if (isNaN(c)) {
                squares.push(
                    new Square(`${x}-${y}`, x, y, (x+y)%2==0, createPieceFromFen(c))
                )
            } else {
                squares.push(...Array.from({length: parseInt(c)}, (_,i) =>
                    new Square(`${x+i}-${y}`, x+i, y, (x+i+y)%2==0)
                ))
            }
        }

        squareRows.push(squares)
    }

    return squareRows
}

class Board {
    constructor(
        player = 'white',
        fen = 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1'
    ) {
        const fenSplit = fen.split(' ')

        if (fenSplit.length != 6) {
            throw new Error('Invalid FEN string')
        }

        this.player = player == 'white' ? 0 : 1
        this.currentPlayer = fenSplit[1] == 'w' ? 0 : 1
        this.players = ['white', 'black']
        this.halfMoveClock = parseInt(fenSplit[4])
        this.fullMoveNumber = parseInt(fenSplit[5])
        this.squares = createSquaresFromFen(fenSplit[0])
    }

    getPieceMoves(x, y) {
        const piece = this.squares[y][x].piece
        return piece ? piece.constructor.getPieceMoves(this.squares, x, y) : []
    }

    movePiece(fromX, fromY, toX, toY) {
        const fromSquare = this.squares[fromY][fromX]
        const toSquare = this.squares[toY][toX]

        if (!fromSquare.piece
            || !this.getPieceMoves(fromX, fromY).some((m) => m.x == toX && m.y == toY)
            || fromSquare.piece.color != this.players[this.currentPlayer]
        ) {
            return false
        }

        toSquare.piece = fromSquare.piece
        fromSquare.piece = null

        this.currentPlayer = (this.currentPlayer+1) % this.players.length
        this.halfMoveClock += 1
        if (this.currentPlayer == 0) {
            this.fullMoveNumber += 1
        }
        return true
    }
}

class Square {
    constructor(id, x, y, isLight, piece = null) {
        this.id = id
        this.x = x
        this.y = y
        this.isLight = isLight
        this.piece = piece
    }
}

class Piece {
    constructor(color, name) {
        this.color = color
        this.name = name
    }

    static getPieceMoves() {
        throw new Error('getPieceMoves() not implemented')
    }

    static getMovesInDirection(squares, x, y, xDir, yDir, max = 8) {
        const moves = []
        const piece = squares[y][x].piece
        const xMin = 0
        const xMax = squares[0].length - 1
        const yMin = 0
        const yMax = squares.length - 1

        let cx = x + xDir
        let cy = y + yDir
        let counter = 0
        while (
            cx >= xMin
            && cx <= xMax
            && cy >= yMin
            && cy <= yMax
            && squares[cy][cx].piece?.color != piece.color
            && counter < max
        ) {
            moves.push({x: cx, y: cy})

            if (squares[cy][cx].piece) {
                break
            }

            cy += yDir
            cx += xDir
            counter++
        }
        return moves
    }
}

class Pawn extends Piece {
    constructor(color, xDir, yDir, moved = false, enPassant = null) {
        super(color, 'p')
        this.xDir = xDir
        this.yDir = yDir
        this.moved = moved
        this.enPassant = enPassant
    }

    static getPieceMoves(squares, x, y) {
        const pawn = squares[y][x].piece
        const moves = []
        const xMin = 0
        const xMax = squares[0].length - 1
        const yMin = 0
        const yMax = squares.length - 1


        const straights = []
        const diagonals = []
        if (pawn.xDir > 0) {
            if (x+1<=xMax) { straights.push({x: x+1, y: y}) }
            if (!pawn.moved && x+2<=xMax) { straights.push({x: x+2, y: y}) }
            if (x+1<=xMax && y+1<=yMax) { diagonals.push({x: x+1, y: y+1, xSide: x, ySide: y+1}) }
            if (x+1<=xMax && y-1>=yMin) { diagonals.push({x: x+1, y: y-1, xSide: x, ySide: y-1}) }
        } else if (pawn.xDir < 0) {
            if (x-1>=xMin) { straights.push({x: x-1, y: y}) }
            if (!pawn.moved && x-2>=xMin) { straights.push({x: x-2, y: y}) }
            if (x-1>=xMin && y+1<=yMax) { diagonals.push({x: x-1, y: y+1, xSide: x, ySide: y+1}) }
            if (x-1>=xMin && y-1>=yMin) { diagonals.push({x: x-1, y: y-1, xSide: x, ySide: y-1}) }
        } else if (pawn.yDir > 0) {
            if (y+1<=yMax) { straights.push({x: x, y: y+1}) }
            if (!pawn.moved && y+2<=yMax) { straights.push({x: x, y: y+2}) }
            if (x+1<=xMax && y+1<=yMax) { diagonals.push({x: x+1, y: y+1, xSide: x+1, ySide: y}) }
            if (x-1>=xMin && y+1<=yMax) { diagonals.push({x: x-1, y: y+1, xSide: x-1, ySide: y}) }
        } else if (pawn.yDir < 0) {
            if (y-1>=yMin) { straights.push({x: x, y: y-1}) }
            if (!pawn.moved && y-2>=yMin) { straights.push({x: x, y: y-2}) }
            if (x+1<=xMax && y-1>=yMin) { diagonals.push({x: x+1, y: y-1, xSide: x+1, ySide: y}) }
            if (x-1>=xMin && y-1>=yMin) { diagonals.push({x: x-1, y: y-1, xSide: x-1, ySide: y}) }
        }

        for (const straight of straights) {
            if (!squares[straight.y][straight.x].piece) {
                moves.push({x: straight.x, y: straight.y})
            } else {
                break
            }
        }

        for (const diagonal of diagonals) {
            if (squares[diagonal.y][diagonal.x].piece
                && squares[diagonal.y][diagonal.x].piece.color != pawn.color) {
                moves.push({x: diagonal.x, y: diagonal.y})
            } else if (squares[diagonal.ySide][diagonal.xSide].piece
                && squares[diagonal.ySide][diagonal.xSide].piece.color != pawn.color
                && squares[diagonal.ySide][diagonal.xSide].piece.name == 'p'
                && squares[diagonal.ySide][diagonal.xSide].piece.enPassant.x == diagonal.x
                && squares[diagonal.ySide][diagonal.xSide].piece.enPassant.y == diagonal.y) {
                moves.push({x: diagonal.x, y: diagonal.y})
            }
        }

        return moves
    }
}

class Knight extends Piece {
    constructor(color) {
        super(color, 'n')
    }

    static getPieceMoves(squares, x, y) {
        return [
            {x: 1, y: 2},
            {x: 1, y: -2},
            {x: -1, y: 2},
            {x: -1, y: -2},
            {x: 2, y: 1},
            {x: 2, y: -1},
            {x: -2, y: 1},
            {x: -2, y: -1},
        ].flatMap(d => Piece.getMovesInDirection(squares, x, y, d.x, d.y, 1))
    }
}

class Bishop extends Piece {
    constructor(color) {
        super(color, 'b')
    }

    static getPieceMoves(squares, x, y) {
        return [
            {x: 1, y: 1},
            {x: 1, y: -1},
            {x: -1, y: 1},
            {x: -1, y: -1},
        ].flatMap(d => Piece.getMovesInDirection(squares, x, y, d.x, d.y))
    }
}

class Rook extends Piece {
    constructor(color, moved = false) {
        super(color, 'r')
        this.moved = moved
    }

    static getPieceMoves(squares, x, y) {
        return [
            {x: 1, y: 0},
            {x: -1, y: 0},
            {x: 0, y: 1},
            {x: 0, y: -1},
        ].flatMap(d => Piece.getMovesInDirection(squares, x, y, d.x, d.y))
    }
}

class Queen extends Piece {
    constructor(color) {
        super(color, 'q')
    }

    static getPieceMoves(squares, x, y) {
        return [
            {x: 1, y: 1},
            {x: 1, y: -1},
            {x: -1, y: 1},
            {x: -1, y: -1},
            {x: 1, y: 0},
            {x: -1, y: 0},
            {x: 0, y: 1},
            {x: 0, y: -1},
        ].flatMap(d => Piece.getMovesInDirection(squares, x, y, d.x, d.y))
    }
}

class King extends Piece {
    constructor(color, moved = false, inCheck = false, inCheckmate = false) {
        super(color, 'k')
        this.moved = moved
        this.inCheck = inCheck
        this.inCheckmate = inCheckmate
    }

    static getPieceMoves(squares, x, y) {
        return [
            {x: 1, y: 1},
            {x: 1, y: -1},
            {x: -1, y: 1},
            {x: -1, y: -1},
            {x: 1, y: 0},
            {x: -1, y: 0},
            {x: 0, y: 1},
            {x: 0, y: -1},
        ].flatMap(d => Piece.getMovesInDirection(squares, x, y, d.x, d.y, 1))
    }
}

export default Board
