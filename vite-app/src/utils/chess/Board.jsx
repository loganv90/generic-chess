import { Square } from '/src/utils/chess/Square'

class Board {
    constructor(
        localPlayer=0,
        currentPlayer=0,
        players=['white', 'black'],
        fen='rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1'
    ) {
        this.localPlayer = localPlayer
        this.currentPlayer = currentPlayer
        this.players = players

        this.enPassant = {}
        for (const player of this.players) {
            this.enPassant[player] = {}
        }

        const fenSplit = fen.split(' ')
        if (fenSplit.length != 6) {
            throw new Error('Invalid FEN string')
        }

        this.halfMoveClock = parseInt(fenSplit[4])
        this.fullMoveNumber = parseInt(fenSplit[5])

        this.squares = this.createSquaresFromFen(fenSplit[0])
        if (this.squares.length < 1 || this.squares[0].length < 1) {
            throw new Error('Invalid FEN string')
        }

        this.xMin = 0
        this.xMax = this.squares[0].length - 1
        this.yMin = 0
        this.yMax = this.squares.length - 1
    }

    createSquaresFromFen(fenPieces) {
        const squareRows = []
        const rows = fenPieces.split('/')
    
        for (const row of rows) {
            const squares = []
    
            for (const c of row) {
                let x = squares.length
                let y = squareRows.length
    
                if (isNaN(c)) {
                    squares.push(
                        new Square(`${x}-${y}`, x, y, (x+y)%2==0, c)
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

    getPiece(x, y) {
        return this.squares[y]?.[x]?.piece
    }

    setPiece(x, y, piece) {
        if (this.squares[y]?.[x]) {
            this.squares[y][x].piece = piece
        }
    }

    getEnPassant(x, y) {
        return this.players
        .map(player => this.enPassant[player])
        .filter(enPassant => enPassant.x === x && enPassant.y === y)
    }

    setEnPassant(x, y, xPiece, yPiece, color) {
        this.enPassant[color] = {x: x, y: y, xPiece: xPiece, yPiece: yPiece}
    }

    clearEnPassant(color) {
        this.enPassant[color] = {}
    }

    getMoves(x, y) {
        return this.getPiece(x, y)?.getMoves()
        .filter(() => this.filterSameColorSource(x, y))
        .flatMap(m => m.options?.isDirection ? this.mapToDirection(x, y, m) : this.mapToDestination(x, y, m))
        .filter(m => this.filterOutOfBounds(m))
        .filter(m => this.filterSameColorDestination(m))
        .map(m => m.options?.canPromote ? this.mapCanPromote(x, y, m) : m)
        .map(m => m.options?.canEnPassant ? this.mapCanEnPassant(m) : m)
        .filter(m => this.filterNoCapture(m))
        .filter(m => this.filterMustCapture(m))
        .filter(m => this.filterMustCross(m))
        .map(m => this.mapRemoveUnusedOptions(m)) ?? []
    }

    mapToDirection(x, y, m) {
        delete m.options?.isDirection
        delete m.options?.enPassant
        delete m.options?.mustCross
        const moves = []
        let cx = x+m.x
        let cy = y+m.y
        while (cx >= this.xMin
            && cx <= this.xMax
            && cy >= this.yMin
            && cy <= this.yMax
        ) {
            moves.push({...m, x: cx, y: cy})
            if (this.getPiece(cx, cy)) {
                break
            }
            cx += m.x
            cy += m.y
        }
        return moves
    }

    mapToDestination(x, y, m) {
        if (m.options?.enPassant) {
            m.options.enPassant.x += x
            m.options.enPassant.y += y
        }
        if (m.options?.mustCross) {
            m.options.mustCross.x += x
            m.options.mustCross.y += y
        }
        m.x += x
        m.y += y
        return m
    }

    filterOutOfBounds(m) {
        return m.x >= this.xMin && m.x <= this.xMax && m.y >= this.yMin && m.y <= this.yMax
    }

    filterSameColorDestination(m) {
        return this.getPiece(m.x, m.y)?.color !== this.players[this.currentPlayer]
    }

    filterSameColorSource(x, y) {
        return this.getPiece(x, y)?.color === this.players[this.currentPlayer]
    }

    filterNoCapture(m) {
        return m.options?.noCapture ? !this.getPiece(m.x, m.y) : true
    }

    filterMustCapture(m) {
        return m.options?.mustCapture ? m.options?.canEnPassant ? true : this.getPiece(m.x, m.y) : true
    }

    filterMustCross(m) {
        return m.options?.mustCross ? !this.getPiece(m.options.mustCross.x, m.options.mustCross.y) : true
    }

    mapCanPromote(x, y, m) {
        if (m.y == this.yMin && y != this.yMin
            || m.y == this.yMax && y != this.yMax
            || m.x == this.xMin && x != this.xMin
            || m.x == this.xMax && x != this.xMax
        ) {
            return m
        }
        delete m.options?.canPromote
        return m
    }

    mapCanEnPassant(m) {
        if (this.players
            .filter(player => player !== this.players[this.currentPlayer])
            .map(player => this.enPassant[player])
            .filter(enPassant => enPassant.x === m.x && enPassant.y === m.y)
            .length > 0
        ) {
            return m
        }
        delete m.options?.canEnPassant
        return m
    }

    mapRemoveUnusedOptions(m) {
        delete m.options?.noCapture
        delete m.options?.mustCapture
        delete m.options?.mustCross
        if (m.options && Object.keys(m.options).length === 0) {
            delete m.options
        }
        return m
    }

    increment() {
        this.currentPlayer = (this.currentPlayer+1) % this.players.length
        this.halfMoveClock += 1
        if (this.currentPlayer == 0) {
            this.fullMoveNumber += 1
        }
    }

    decrement() {
        this.currentPlayer = this.currentPlayer-1 < 0 ? this.players.length-1 : this.currentPlayer-1
        this.halfMoveClock -= 1
        if (this.currentPlayer == this.players.length-1) {
            this.fullMoveNumber -= 1
        }
    }
}

export { Board }
