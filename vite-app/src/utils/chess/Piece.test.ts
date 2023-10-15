import { describe, it, expect, vi } from 'vitest'
import { Piece, Pawn, Knight, Bishop, Rook, Queen, King } from './Piece'
import { Board } from './Board'
import { SimpleMove } from './Move'

vi.mock('./Move', () => {
    const SimpleMove = vi.fn()
    return { SimpleMove }
})

vi.mock('./Board', () => {
    const Board = vi.fn(() => ({
        xMin: 0,
        xMax: 7,
        yMin: 0,
        yMax: 7,
        getPiece: vi.fn(() => null),
        getEnPassants: vi.fn(() => []),
    }))
    return { Board }
})

describe('Pawn', () => {
    describe('getMoves', () => {
        it('should return correct moves when unmoved', () => {
            const board = new Board()
            const pawn = new Pawn('white', 0, 1)

            const moves = pawn.getMoves(3, 3, board)

            expect(moves.length).toBe(2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 5)
        })

        it('should return correct moves when moved', () => {
            const board = new Board()
            const pawn = new Pawn('white', 0, 1, true)

            const moves = pawn.getMoves(3, 3, board)

            expect(moves.length).toBe(1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 4)
        })

        it('should return correct moves when capturing', () => {
            const board = new Board()
            const pawn = new Pawn('white', 0, 1)
            const blackPawn = new Pawn('black', 0, -1)
            board.getPiece = vi.fn((x: number, y: number): Piece | null => {
                if (x === 4 && y === 4) {
                    return blackPawn
                }
                return null
            })

            const moves = pawn.getMoves(3, 3, board)

            expect(moves.length).toBe(3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 5)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 4, 4)
        })

        it('should return correct moves when capturing en passant', () => {
            const board = new Board()
            const pawn = new Pawn('white', 0, 1)
            board.getEnPassants = vi.fn(() => [{
                xTarget: 4,
                yTarget: 4,
                xPiece: 4,
                yPiece: 3
            }])

            const moves = pawn.getMoves(3, 3, board)

            expect(moves.length).toBe(4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 5)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 4, 4)
        })
    })
})

describe('Knight', () => {
    describe('getMoves', () => {
        it('should return knight moves', () => {
            const board = new Board()
            const knight = new Knight('white')

            const moves = knight.getMoves(1, 1, board)

            expect(moves.length).toBe(4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 2, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 3, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 0, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 3, 0)
        })
    })
})

describe('Bishop', () => {
    describe('getMoves', () => {
        it('should return bishop moves', () => {
            const board = new Board()
            const bishop = new Bishop('white')

            const moves = bishop.getMoves(1, 1, board)

            expect(moves.length).toBe(9)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 0, 0)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 0, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 2, 0)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 2, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 3, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 4, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 5, 5)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 6, 6)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 7, 7)
        })
    })
})

describe('Rook', () => {
    describe('getMoves', () => {
        it('should return rook moves', () => {
            const board = new Board()
            const rook = new Rook('white')

            const moves = rook.getMoves(1, 1, board)

            expect(moves.length).toBe(14)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 0, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 2, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 3, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 4, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 5, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 6, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 7, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 0)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 5)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 6)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 7)
        })
    })
})

describe('Queen', () => {
    describe('getMoves', () => {
        it('should return queen moves', () => {
            const board = new Board()
            const queen = new Queen('white')

            const moves = queen.getMoves(1, 1, board)

            expect(moves.length).toBe(23)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 0, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 2, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 3, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 4, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 5, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 6, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 7, 1)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 0)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 5)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 6)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 1, 7)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 0, 0)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 0, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 2, 0)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 2, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 3, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 4, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 5, 5)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 6, 6)
            expect(SimpleMove).toHaveBeenCalledWith(board, 1, 1, 7, 7)
        })
    })
})

describe('King', () => {
    describe('getMoves', () => {
        it('should return correct moves when can castle and unmoved', () => {
            const board = new Board()
            const king = new King('white', 0, 1)
            const rook = new Rook('white')
            board.getPiece = vi.fn((x: number, y: number): Piece | null => {
                if (x === 0 && y === 3) {
                    return rook
                } else if (x === 7 && y === 3) {
                    return rook
                }
                return null
            })

            const moves = king.getMoves(3, 3, board)

            expect(moves.length).toBe(10)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 2, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 4, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 2, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 4, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 2, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 4, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 0, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 7, 3)
        })

        it('should return correct moves when can castle and moved', () => {
            const board = new Board()
            const king = new King('white', 0, 1)
            const rook = new Rook('white', true)
            board.getPiece = vi.fn((x: number, y: number): Piece | null => {
                if (x === 0 && y === 3) {
                    return rook
                } else if (x === 7 && y === 3) {
                    return rook
                }
                return null
            })

            const moves = king.getMoves(3, 3, board)

            expect(moves.length).toBe(8)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 2, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 4, 2)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 2, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 4, 3)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 2, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 3, 4)
            expect(SimpleMove).toHaveBeenCalledWith(board, 3, 3, 4, 4)
        })
    })
})