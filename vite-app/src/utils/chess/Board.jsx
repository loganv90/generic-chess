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

        this.minX = 0
        this.maxX = this.squares[0].length - 1
        this.minY = 0
        this.maxY = this.squares.length - 1
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
        return this.squares[y]?.[x]?.getPiece()
    }

    getMoves(x, y) {
        return this.getPiece(x, y).getMoves()
        .flatMap(m => m.options?.isDirection ? this.mapMoveDirection(x,y,m) : this.mapMoveDestination(x,y,m))
        .filter(m => this.getPiece(m.x, m.y)?.getColor() !== this.players[this.currentPlayer])
        .map(m => m.options?.canPromote ? this.mapMoveCanPromote(x,y,m) : m)
        .map(m => m.options?.canEnPassant ? this.mapMoveCanEnPassant(m) : m)
        .filter(m => m.options?.noCapture ? !this.getPiece(m.x, m.y) : true)
        .map(m => {delete m.options.noCapture; return m})
        .filter(m => m.options?.mustCapture ? m.options?.canEnPassant ? true : this.getPiece(m.x, m.y) : true)
        .map(m => {delete m.options.mustCapture; return m})
    }

    mapMoveDirection(x, y, m) {
        const moves = []
        let cx = x+m.x
        let cy = y+m.y
        do {
            delete m.options.isDirection
            moves.push({...m, x: cx, y: cy})
            cy += m.x
            cx += m.y
        } while (
            cx >= this.xMin
            && cx <= this.xMax
            && cy >= this.yMin
            && cy <= this.yMax
            && this.getPiece()
        )
        return moves
    }

    mapMoveDestination(x, y, m) {
        if (m.options?.enPassant) {
            m.options.enPassant.x = x+m.options.enPassant.x
            m.options.enPassant.y = y+m.options.enPassant.y
        }
        return {...m, x: x+m.x, y: y+m.y}
    }

    mapMoveCanPromote(x, y, m) {
        if (m.y == this.minY && y != this.minY
            || m.y == this.maxY && y != this.maxY
            || m.x == this.minX && x != this.minX
            || m.x == this.maxX && x != this.maxX
        ) {
            m.options.canPromote = true
        } else {
            delete m.options.canPromote
        }
        return m
    }

    mapMoveCanEnPassant(m) {
        const canEnPassant = [{x: 1, y:0}, {x: -1, y:0}, {x: 0, y:1}, {x: 0, y:-1}].map(d => {
            const enPassant = this.getPiece(m.x+d.x, m.y+d.y)?.getEnPassant()
            if (enPassant && enPassant.x === m.x && enPassant.y === m.y) {
                return {x: m.x+d.x, y: m.y+d.y}
            }
        }).filter(e => e)
        if (canEnPassant.length > 0) {
            m.options.canEnPassant = canEnPassant
        } else {
            delete m.options.canEnPassant
        }
        return m
    }

    // getPawnMoves(x, y) {
    //     const pawn = this.squares[y][x].piece
    //     const moves = []
    //     const xMin = 0
    //     const xMax = this.squares[0].length - 1
    //     const yMin = 0
    //     const yMax = this.squares.length - 1


    //     const straights = []
    //     const diagonals = []
    //     if (pawn.xDir > 0) {
    //         if (x+1<=xMax) { straights.push({x: x+1, y: y}) }
    //         if (!pawn.moved && x+2<=xMax) { straights.push({x: x+2, y: y}) }
    //         if (x+1<=xMax && y+1<=yMax) { diagonals.push({x: x+1, y: y+1, xSide: x, ySide: y+1}) }
    //         if (x+1<=xMax && y-1>=yMin) { diagonals.push({x: x+1, y: y-1, xSide: x, ySide: y-1}) }
    //     } else if (pawn.xDir < 0) {
    //         if (x-1>=xMin) { straights.push({x: x-1, y: y}) }
    //         if (!pawn.moved && x-2>=xMin) { straights.push({x: x-2, y: y}) }
    //         if (x-1>=xMin && y+1<=yMax) { diagonals.push({x: x-1, y: y+1, xSide: x, ySide: y+1}) }
    //         if (x-1>=xMin && y-1>=yMin) { diagonals.push({x: x-1, y: y-1, xSide: x, ySide: y-1}) }
    //     } else if (pawn.yDir > 0) {
    //         if (y+1<=yMax) { straights.push({x: x, y: y+1}) }
    //         if (!pawn.moved && y+2<=yMax) { straights.push({x: x, y: y+2}) }
    //         if (x+1<=xMax && y+1<=yMax) { diagonals.push({x: x+1, y: y+1, xSide: x+1, ySide: y}) }
    //         if (x-1>=xMin && y+1<=yMax) { diagonals.push({x: x-1, y: y+1, xSide: x-1, ySide: y}) }
    //     } else if (pawn.yDir < 0) {
    //         if (y-1>=yMin) { straights.push({x: x, y: y-1}) }
    //         if (!pawn.moved && y-2>=yMin) { straights.push({x: x, y: y-2}) }
    //         if (x+1<=xMax && y-1>=yMin) { diagonals.push({x: x+1, y: y-1, xSide: x+1, ySide: y}) }
    //         if (x-1>=xMin && y-1>=yMin) { diagonals.push({x: x-1, y: y-1, xSide: x-1, ySide: y}) }
    //     }

    //     for (const straight of straights) {
    //         if (!this.squares[straight.y][straight.x].piece) {
    //             moves.push({x: straight.x, y: straight.y})
    //         } else {
    //             break
    //         }
    //     }

    //     for (const diagonal of diagonals) {
    //         if (this.squares[diagonal.y][diagonal.x].piece
    //             && this.squares[diagonal.y][diagonal.x].piece.color != pawn.color) {
    //             moves.push({x: diagonal.x, y: diagonal.y})
    //         } else if (this.squares[diagonal.ySide][diagonal.xSide].piece
    //             && this.squares[diagonal.ySide][diagonal.xSide].piece.color != pawn.color
    //             && this.squares[diagonal.ySide][diagonal.xSide].piece.name == 'p'
    //             && this.squares[diagonal.ySide][diagonal.xSide].piece.enPassant.x == diagonal.x
    //             && this.squares[diagonal.ySide][diagonal.xSide].piece.enPassant.y == diagonal.y) {
    //             moves.push({x: diagonal.x, y: diagonal.y})
    //         }
    //     }

    //     return moves
    // }

    // static move(squares, fromX, fromY, toX, toY) {
    //     const pawn = squares[fromY][fromX].piece
    //     const dx = toX - fromX
    //     const dy = toY - fromY

    //     if (Math.abs(dx) + Math.abs(dy) > 1) {
    //         pawn.enPassant = {
    //             x: dx>0 ? fromX+1 : dx<0 ? fromX-1 : fromX,
    //             y: dy>0 ? fromY+1 : dy<0 ? fromY-1 : fromY,
    //         }
    //     }

    //     return super.move(squares, fromX, fromY, toX, toY)
    // }

    // movePiece(fromX, fromY, toX, toY) {
    //     const piece = this.squares[fromY][fromX].piece

    //     if (!piece
    //         || piece.color != this.players[this.currentPlayer]
    //         || !this.getPieceMoves(fromX, fromY).some((m) => m.x == toX && m.y == toY)
    //         || !piece.constructor.move(this.squares, fromX, fromY, toX, toY)
    //     ) {
    //         return false
    //     }

    //     this.currentPlayer = (this.currentPlayer+1) % this.players.length
    //     this.halfMoveClock += 1
    //     if (this.currentPlayer == 0) {
    //         this.fullMoveNumber += 1

    //     }

    //     this.squares.forEach((row) => row.forEach((s) => {
    //         if (s.piece
    //             && s.piece.color === this.players[this.currentPlayer]
    //             && s.piece.enPassant
    //         ) {
    //             s.piece.enPassant = null
    //         }
    //     }))

    //     return true
    // }
}

export { Board }
