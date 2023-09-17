import { PieceMove, BoardMove, EnPassantMap, EnPassant } from './Types'
import { Square } from './Square'
import { Piece } from './Piece'

class Board {
    readonly squares: Square[][]
    private localPlayer: number
    private currentPlayer: number
    private players: string[]
    private enPassant: EnPassantMap
    private halfMoveClock: number
    private fullMoveNumber: number
    private xMin: number
    private xMax: number
    private yMin: number
    private yMax: number

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

    getLocalPlayer(): number {
        return this.localPlayer
    }

    getCurrentPlayer(): number {
        return this.currentPlayer
    }

    getPiece(x: number, y: number): Piece | null {
        return this.squares[y]?.[x]?.getPiece() ?? null
    }

    setPiece(x: number, y: number, piece: Piece | null): void {
        this.squares[y]?.[x]?.setPiece(piece)
    }

    getEnPassant(color: string): EnPassant | null {
        return this.enPassant[color] ?? null
    }

    setEnPassant(color: string, enPassant: EnPassant | null) {
        this.enPassant[color] = enPassant
    }

    getEnPassants(color: string, x: number, y: number): EnPassant[] {
        return Object.entries(this.enPassant)
        .filter((e): e is [string, EnPassant] => e[0] !== color)
        .filter((e): e is [string, EnPassant] => e[1]?.x === x && e[1].y === y)
        .map(e => e[1])
    }

    getMoves(x: number, y: number): BoardMove[] {
        return this.getPiece(x, y)?.getMoves()
        .filter(() => this.filterSameColorSource(x, y))
        .flatMap(m => m.options.direction ? this.mapToDirection(x, y, m) : this.mapToDestination(x, y, m))
        .filter(m => this.filterOutOfBounds(m))
        .filter(m => this.filterSameColorDestination(m))
        .map(m => m.options.canPromote ? this.mapCanPromote(x, y, m) : m)
        .map(m => m.options.canCaptureEnPassant ? this.mapCanCaptureEnPassant(m) : m)
        .filter(m => this.filterNoCapture(m))
        .filter(m => this.filterMustCapture(m))
        .filter(m => this.filterMustRevealEnPassant(m))
        .map(m => this.mapToBoardMove(x, y, m)) ?? []
    }

    increment(): void {
        this.currentPlayer = (this.currentPlayer+1) % this.players.length
        this.halfMoveClock += 1
        if (this.currentPlayer == 0) {
            this.fullMoveNumber += 1
        }
    }

    decrement(): void {
        this.currentPlayer = this.currentPlayer-1 < 0 ? this.players.length-1 : this.currentPlayer-1
        this.halfMoveClock -= 1
        if (this.currentPlayer == this.players.length-1) {
            this.fullMoveNumber -= 1
        }
    }

    private createSquaresFromFen(fenPieces: string): Square[][] {
        const squareRows = []
        const rows = fenPieces.split('/')
    
        for (const row of rows) {
            const squares = []
    
            for (const c of row) {
                const x: number = squares.length
                const y: number = squareRows.length
                const n: number = parseInt(c)
    
                if (isNaN(n)) {
                    squares.push(
                        new Square(`${x}-${y}`, x, y, c)
                    )
                } else {
                    squares.push(...Array.from({length: n}, (_,i) =>
                        new Square(`${x+i}-${y}`, x+i, y)
                    ))
                }
            }
    
            squareRows.push(squares)
        }
    
        return squareRows
    }

    private mapToDirection(x: number, y: number, m: PieceMove): PieceMove[] {
        delete m.options.direction
        delete m.options.mustRevealEnPassant
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

    private mapToDestination(x: number, y: number, m: PieceMove): PieceMove {
        if (m.options.mustRevealEnPassant) {
            m.options.mustRevealEnPassant.x += x
            m.options.mustRevealEnPassant.y += y
        }
        m.x += x
        m.y += y
        return m
    }

    private filterOutOfBounds(m: PieceMove): boolean {
        return m.x >= this.xMin && m.x <= this.xMax && m.y >= this.yMin && m.y <= this.yMax
    }

    private filterSameColorDestination(m: PieceMove): boolean {
        return this.getPiece(m.x, m.y)?.color !== this.players[this.currentPlayer]
    }

    private filterSameColorSource(x: number, y: number): boolean {
        return this.getPiece(x, y)?.color === this.players[this.currentPlayer]
    }

    private filterNoCapture(m: PieceMove): boolean {
        return m.options.noCapture ? !this.getPiece(m.x, m.y) : true
    }

    private filterMustCapture(m: PieceMove): boolean {
        return Boolean(m.options.mustCapture ? m.options.canCaptureEnPassant ? true : this.getPiece(m.x, m.y) : true)
    }

    private filterMustRevealEnPassant(m: PieceMove): boolean {
        return m.options.mustRevealEnPassant ? !this.getPiece(m.options.mustRevealEnPassant.x, m.options.mustRevealEnPassant.y) : true
    }

    private mapCanPromote(x: number, y: number, m: PieceMove): PieceMove {
        if (m.y == this.yMin && y != this.yMin
            || m.y == this.yMax && y != this.yMax
            || m.x == this.xMin && x != this.xMin
            || m.x == this.xMax && x != this.xMax
        ) {
            return m
        }
        delete m.options.canPromote
        return m
    }

    private mapCanCaptureEnPassant(m: PieceMove): PieceMove {
        if (this.getEnPassants(this.players[this.currentPlayer], m.x, m.y).length > 0) {
            return m
        }
        delete m.options.canCaptureEnPassant
        return m
    }

    private mapToBoardMove(x: number, y: number, m: PieceMove): BoardMove {
        const boardMove: BoardMove = {xFrom: x, yFrom: y, xTo: m.x, yTo: m.y, options: {}}
        if (m.options.canCaptureEnPassant) {
            boardMove.options.captureEnPassant = m.options.canCaptureEnPassant
        }
        if (m.options.mustRevealEnPassant) {
            boardMove.options.revealEnPassant = m.options.mustRevealEnPassant
        }
        if (m.options.canPromote) {
            boardMove.options.promote = true
        }
        return boardMove
    }
}

export { Board }