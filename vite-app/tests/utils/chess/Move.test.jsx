import { describe, it, expect, vi } from 'vitest'
import { SimpleMove } from '/src/utils/chess/Move'
import { Board } from '/src/utils/chess/Board'
import { Piece } from '/src/utils/chess/Piece'

vi.mock('/src/utils/chess/Piece', () => {
    const Piece = vi.fn(() => ({
        move: vi.fn(),
    }))
    return { Piece }
})

vi.mock('/src/utils/chess/Board', () => {
    const Board = vi.fn(() => ({
        setPiece: vi.fn(),
        increment: vi.fn(),
        decrement: vi.fn(),
    }))
    return { Board }
})

describe('SimpleMove', () => {
    it('should get and set the correct pieces', () => {
        const piece1 = new Piece()
        const piece2 = new Piece()
        const piece3 = new Piece()
        const board = new Board()

        piece1.copy = vi.fn(() => piece3)
        board.getPiece = vi.fn((x, y) => {
            if (x == 0 && y == 0) {
                return piece1
            }
            if (x == 1 && y == 1) {
                return piece2
            }
        })

        const move = new SimpleMove(board, 0, 0, 1, 1)
        expect(board.getPiece).toBeCalledWith(0, 0)
        expect(board.getPiece).toBeCalledWith(1, 1)

        expect(move.execute()).toBe(true)
        expect(board.setPiece).toBeCalledWith(0, 0, null)
        expect(board.setPiece).toBeCalledWith(1, 1, piece3)

        expect(move.undo()).toBe(true)
        expect(board.setPiece).toBeCalledWith(0, 0, piece1)
        expect(board.setPiece).toBeCalledWith(1, 1, piece2)
    })
})