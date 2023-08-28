import { describe, it, expect, vi } from 'vitest'
import { SimpleMove, RevealEnPassantMove, CaptureEnPassantMove } from '/src/utils/chess/Move'
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
        setEnPassant: vi.fn(),
    }))
    return { Board }
})

describe('SimpleMove', () => {
    it('should get and set the correct pieces', () => {
        const board = new Board()
        const pieceWhite = new Piece()
        const pieceBlack = new Piece()
        const pieceWhiteNew = new Piece()

        const enPassantWhite = {x: 1, y: 1, xPiece: 2, yPiece: 2}
        const enPassantBlack = {x: 1, y: 1, xPiece: 3, yPiece: 3}
        pieceWhite.copy = vi.fn(() => pieceWhiteNew)
        pieceWhite.color = 'white'
        pieceBlack.color = 'black'
        board.getPiece = vi.fn((x, y) => {
            if (x == 0 && y == 0) {
                return pieceWhite
            }
            if (x == 1 && y == 1) {
                return pieceBlack
            }
        })
        board.getEnPassant = vi.fn((color) => {
            if (color == 'white') {
                return enPassantWhite
            }
            if (color == 'black') {
                return enPassantBlack
            }
        })
        board.enPassant = {white: enPassantWhite, black: enPassantBlack}

        const move = new SimpleMove(board, 0, 0, 1, 1)
        expect(board.getPiece).toBeCalledWith(0, 0)
        expect(board.getPiece).toBeCalledWith(1, 1)

        expect(move.execute()).toBe(true)
        expect(board.setPiece).toBeCalledWith(0, 0, null)
        expect(board.setPiece).toBeCalledWith(1, 1, pieceWhiteNew)
        expect(board.increment).toBeCalled()
        expect(board.setEnPassant).toBeCalledWith(pieceWhite.color, null)

        expect(move.undo()).toBe(true)
        expect(board.setPiece).toBeCalledWith(0, 0, pieceWhite)
        expect(board.setPiece).toBeCalledWith(1, 1, pieceBlack)
        expect(board.decrement).toBeCalled()
        expect(board.setEnPassant).toBeCalledWith(pieceWhite.color, enPassantWhite)
    })
})

describe('RevealEnPassantMove', () => {
    it('should get and set the correct pieces', () => {
        const board = new Board()
        const piece = new Piece()

        const enPassantWhite = {x: 2, y: 2, xPiece: 1, yPiece: 1}
        piece.color = 'white'
        piece.copy = vi.fn(() => piece)
        board.getPiece = vi.fn(() => piece)
        board.getEnPassant = vi.fn()

        const move = new RevealEnPassantMove(board, 0, 0, 1, 1, 2, 2)

        expect(move.execute()).toBe(true)
        expect(board.setEnPassant.mock.calls.length).toBe(2)
        expect(board.setEnPassant.mock.calls.at(-1)).toEqual([piece.color, enPassantWhite])

        expect(move.undo()).toBe(true)

    })
})

describe('CaptureEnPassantMove', () => {
    it('should get and set the correct pieces', () => {
        const board = new Board()
        const pieceWhite = new Piece()
        const pieceBlack = new Piece()
        const pieceGrey = new Piece()

        const enPassantGrey = {x: 2, y: 2, xPiece: 1, yPiece: 1}
        const enPassantBlack = {x: 2, y: 2, xPiece: 3, yPiece: 3}
        pieceWhite.color = 'white'
        pieceWhite.copy = vi.fn(() => pieceWhite)
        pieceBlack.copy = vi.fn(() => pieceBlack)
        pieceGrey.copy = vi.fn(() => pieceGrey)
        board.getPiece = vi.fn((x, y) => {
            if (x == 1 && y == 1) {
                return pieceGrey
            }
            if (x == 3 && y == 3) {
                return pieceBlack
            }
            return pieceWhite
        })
        board.getEnPassant = vi.fn()
        board.getEnPassants = vi.fn((color, x, y) => {
            if (color === 'white' && x === 1 && y === 1) {
                return [enPassantGrey, enPassantBlack]
            }
        })

        const move = new CaptureEnPassantMove(board, 0, 0, 1, 1)

        expect(move.execute()).toBe(true)
        expect(board.setPiece).toBeCalledWith(1, 1, null)
        expect(board.setPiece).toBeCalledWith(3, 3, null)

        expect(move.undo()).toBe(true)
        expect(board.setPiece).toBeCalledWith(1, 1, pieceGrey)
        expect(board.setPiece).toBeCalledWith(3, 3, pieceBlack)
    })
})