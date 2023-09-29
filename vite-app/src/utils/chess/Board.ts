import { Square } from './Square'
import { Piece } from './Piece'
import { Move } from './Move'

type EnPassantMap = {
    [key: string]: EnPassant | null,
}

type EnPassant = {
    xTarget: number,
    yTarget: number,
    xPiece: number,
    yPiece: number,
}

class Board {
    readonly squares: Square[][]
    readonly xMin: number
    readonly xMax: number
    readonly yMin: number
    readonly yMax: number
    private localPlayer: number
    private currentPlayer: number
    private players: string[]
    private enPassant: EnPassantMap
    private halfMoveClock: number
    private fullMoveNumber: number

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
        .filter((e): e is [string, EnPassant] => e[1]?.xTarget === x && e[1].yTarget === y)
        .map(e => e[1])
    }

    getMoves(x: number, y: number): Move[] {
        const piece = this.getPiece(x, y)

        if (piece && piece.color === this.players[this.currentPlayer]) {
            return piece.getMoves(x, y, this)
        }

        return []
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
}

export { Board }