import { Board } from './Board'
import { Move, SimpleMove } from './Move'

abstract class Piece {
    readonly color: string
    readonly moved: boolean

    constructor(color: string, moved=false) {
        this.color = color
        this.moved = moved
    }

    getMoves(xFrom: number, yFrom: number, board: Board): Move[] {
        const moves: Move[] = []

        this.addMoves(xFrom, yFrom, board, moves)

        return moves
    }

    abstract copy(): Piece

    protected abstract addMoves(xFrom: number, yFrom: number, board: Board, moves: Move[]): void

    protected static movedOutOfBounds = (
        xTo: number,
        yTo: number,
        board: Board
    ): boolean => {
        return xTo < board.xMin
            || xTo > board.xMax
            || yTo < board.yMin
            || yTo > board.yMax
    }
    
    protected static movedToPromotionSquare = (
        xTo: number,
        yTo: number,
        xFrom: number,
        yFrom: number,
        board: Board
    ): boolean => {
        return (xTo === board.xMin && xFrom !== board.xMin)
            || (xTo === board.xMax && xFrom !== board.xMax)
            || (yTo === board.yMin && yFrom !== board.yMin)
            || (yTo === board.yMax && yFrom !== board.yMax)
    }

    protected static addDirection = (
        xFrom: number,
        yFrom: number,
        board: Board,
        moves: Move[],
        dx: number,
        dy: number,
        color: string,
    ): void => {
        let cx = xFrom
        let cy = yFrom

        while (!Piece.movedOutOfBounds(cx, cy, board)) {
            cx += dx
            cy += dy
            const piece = board.getPiece(cx, cy)

            if (piece && piece.color === color) {
                break
            } else if (piece && piece.color !== color) {
                const simpleMove = new SimpleMove(board, xFrom, yFrom, cx, cy)
                moves.push(simpleMove)
                break
            } else {
                const simpleMove = new SimpleMove(board, xFrom, yFrom, cx, cy)
                moves.push(simpleMove)
            }
        }
    }

    protected static addSimple = (
        xFrom: number,
        yFrom: number,
        board: Board,
        moves: Move[],
        dx: number,
        dy: number,
        color: string,
    ): void => {
        const xTo = xFrom + dx
        const yTo = yFrom + dy
        const piece = board.getPiece(xTo, yTo)

        if (
            !Piece.movedOutOfBounds(xTo, yTo, board)
            && !(piece && piece.color === color)
        ) {
            const simpleMove = new SimpleMove(board, xFrom, yFrom, xTo, yTo)
            moves.push(simpleMove)
        }
    }
}

class Pawn extends Piece {
    private readonly xDir: 1|0|-1
    private readonly yDir: 1|0|-1
    private readonly forward1: {dx: number, dy: number}
    private readonly forward2: {dx: number, dy: number}
    private readonly captures: {dx: number, dy: number}[]

    constructor(color: string, xDir: 1|0|-1, yDir: 1|0|-1, moved=false) {
        super(color, moved)

        if (xDir && yDir) {
            throw new Error('Pawn can not move diagonally')
        } else if (!xDir && !yDir) {
            throw new Error('Pawn must have a direction')
        }

        if (xDir) {
            this.forward1 = {dx: xDir, dy: 0}
            this.forward2 = {dx: xDir*2, dy: 0}
            this.captures = [
                {dx: xDir, dy: 1},
                {dx: xDir, dy: -1},
            ]
        } else {
            this.forward1 = {dx: 0, dy: yDir}
            this.forward2 = {dx: 0, dy: yDir*2}
            this.captures = [
                {dx: 1, dy: yDir},
                {dx: -1, dy: yDir},
            ]
        }

        this.xDir = xDir
        this.yDir = yDir
    }

    copy(): Pawn {
        return new Pawn(this.color, this.xDir, this.yDir, true)
    }

    protected addMoves(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        this.addForward(xFrom, yFrom, board, moves)
        this.addCaptures(xFrom, yFrom, board, moves)
    }

    private addForward(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        const xTo1 = xFrom + this.forward1.dx
        const yTo1 = yFrom + this.forward1.dy

        if (Piece.movedOutOfBounds(xTo1, yTo1, board) || board.getPiece(xTo1, yTo1)) {
            return
        } else if (Piece.movedToPromotionSquare(xTo1, yTo1, xFrom, yFrom, board)) {
            const promotionMove = new SimpleMove(board, xFrom, yFrom, xTo1, yTo1)
            moves.push(promotionMove)
        } else {
            const simpleMove = new SimpleMove(board, xFrom, yFrom, xTo1, yTo1)
            moves.push(simpleMove)
        }

        if (this.moved) {
            return
        }

        const xTo2 = xFrom + this.forward2.dx
        const yTo2 = yFrom + this.forward2.dy

        if (Piece.movedOutOfBounds(xTo2, yTo2, board) || board.getPiece(xTo2, yTo2)) {
            return
        } else if (Piece.movedToPromotionSquare(xTo2, yTo2, xFrom, yFrom, board)) {
            const promotionMove = new SimpleMove(board, xFrom, yFrom, xTo2, yTo2)
            moves.push(promotionMove)
        } else {
            const revealEnPassantMove = new SimpleMove(board, xFrom, yFrom, xTo2, yTo2)
            moves.push(revealEnPassantMove)
        }
    }

    private addCaptures(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        this.captures.forEach(m => {
            const xTo = xFrom + m.dx
            const yTo = yFrom + m.dy
            const piece = board.getPiece(xTo, yTo)
            const enPassants = board.getEnPassants(this.color, xTo, yTo)
            
            if (Piece.movedToPromotionSquare(xTo, yTo, xFrom, yFrom, board)) {
                const promotionMove = new SimpleMove(board, xFrom, yFrom, xTo, yTo)
                moves.push(promotionMove)
            } else if (enPassants.length > 0) {
                const captureEnPassantMove = new SimpleMove(board, xFrom, yFrom, xTo, yTo)
                moves.push(captureEnPassantMove)
            } else if (piece && piece.color !== this.color) {
                const simpleMove = new SimpleMove(board, xFrom, yFrom, xTo, yTo)
                moves.push(simpleMove)
            }
        })
    }
}

class Knight extends Piece {
    private static readonly simples: {dx: number, dy: number}[] = [
        {dx: 1, dy: 2},
        {dx: 1, dy: -2},
        {dx: -1, dy: -2},
        {dx: -1, dy: 2},
        {dx: 2, dy: 1},
        {dx: 2, dy: -1},
        {dx: -2, dy: 1},
        {dx: -2, dy: -1}, 
    ]

    constructor(color: string) {
        super(color)
    }

    copy(): Knight {
        return new Knight(this.color)
    }

    protected addMoves(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        this.addSimples(xFrom, yFrom, board, moves)
    }

    private addSimples(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        Knight.simples.forEach(m => Piece.addSimple(xFrom, yFrom, board, moves, m.dx, m.dy, this.color))
    }
}

class Bishop extends Piece {
    private static readonly directions: {dx: number, dy: number}[] = [
        {dx: 1, dy: 1},
        {dx: 1, dy: -1},
        {dx: -1, dy: 1},
        {dx: -1, dy: -1},
    ]

    constructor(color: string) {
        super(color)
    }

    copy(): Bishop {
        return new Bishop(this.color)
    }

    protected addMoves(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        this.addDirections(xFrom, yFrom, board, moves)
    }

    private addDirections(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        Bishop.directions.forEach(m => Piece.addDirection(xFrom, yFrom, board, moves, m.dx, m.dy, this.color))
    }
}

class Rook extends Piece {
    private static readonly directions: {dx: number, dy: number}[] = [
        {dx: 1, dy: 0},
        {dx: -1, dy: 0},
        {dx: 0, dy: 1},
        {dx: 0, dy: -1},
    ]

    constructor(color: string, moved=false) {
        super(color, moved)
    }

    copy(): Rook {
        return new Rook(this.color, true)
    }

    protected addMoves(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        this.addDirections(xFrom, yFrom, board, moves)
    }
    
    private addDirections(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        Rook.directions.forEach(m => Piece.addDirection(xFrom, yFrom, board, moves, m.dx, m.dy, this.color))
    }
}

class Queen extends Piece {
    private static readonly directions: {dx: number, dy: number}[] = [
        {dx: 1, dy: 1},
        {dx: 1, dy: -1},
        {dx: -1, dy: 1},
        {dx: -1, dy: -1},
        {dx: 1, dy: 0},
        {dx: -1, dy: 0},
        {dx: 0, dy: 1},
        {dx: 0, dy: -1},
    ]

    constructor(color: string) {
        super(color)
    }

    copy(): Queen {
        return new Queen(this.color)
    }

    protected addMoves(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        this.addDirections(xFrom, yFrom, board, moves)
    }
    
    private addDirections(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        Queen.directions.forEach(m => Piece.addDirection(xFrom, yFrom, board, moves, m.dx, m.dy, this.color))
    }
}

class King extends Piece {
    private readonly xDir: 1|0|-1;
    private readonly yDir: 1|0|-1;
    private readonly castles: {dx: number, dy: number}[] = []
    private static readonly simples: {dx: number, dy: number}[] = [
        {dx: 1, dy: 1},
        {dx: 1, dy: -1},
        {dx: -1, dy: 1},
        {dx: -1, dy: -1},
        {dx: 1, dy: 0},
        {dx: -1, dy: 0},
        {dx: 0, dy: 1},
        {dx: 0, dy: -1},
    ]

    constructor(color: string, xDir: 1|0|-1, yDir: 1|0|-1, moved=false) {
        super(color, moved)

        if (xDir && yDir) {
            throw new Error('King can not move diagonally')
        } else if (!xDir && !yDir) {
            throw new Error('King must have a direction')
        }

        if (moved) {
            this.castles = []
        } else if (xDir) {
            this.castles = [
                {dx: 0, dy: 1},
                {dx: 0, dy: -1},
            ]
        } else {
            this.castles = [
                {dx: 1, dy: 0},
                {dx: -1, dy: 0},
            ]
        }

        this.xDir = xDir
        this.yDir = yDir
    }

    copy(): King {
        return new King(this.color, this.xDir, this.yDir, true)
    }

    protected addMoves(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        this.addSimples(xFrom, yFrom, board, moves)
        this.addCastles(xFrom, yFrom, board, moves)
    }
    
    private addSimples(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        King.simples.forEach(m => Piece.addSimple(xFrom, yFrom, board, moves, m.dx, m.dy, this.color))
    }

    private addCastles(xFrom: number, yFrom: number, board: Board, moves: Move[]): void {
        this.castles.forEach(m => {
            let cx = xFrom
            let cy = yFrom

            while (!Piece.movedOutOfBounds(cx, cy, board)) {
                cx += m.dx
                cy += m.dy
                const piece = board.getPiece(cx, cy)

                if (!piece) {
                    continue
                } else if (
                    piece instanceof Rook
                    && piece.color === this.color
                    && !piece.moved
                ) {
                    const castleMove = new SimpleMove(board, xFrom, yFrom, cx, cy)
                    moves.push(castleMove)
                } else {
                    break
                }
            }
        })
    }
}

export { Piece, Pawn, Knight, Bishop, Rook, Queen, King }