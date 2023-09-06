import { describe, it, expect, vi } from 'vitest'
import { SimpleMove, RevealEnPassantMove, CaptureEnPassantMove } from '../../../src/utils/chess/Move'
import { Board } from '../../../src/utils/chess/Board'
import { Pawn } from '../../../src/utils/chess/Piece'
import { BoardMove } from '../../../src/utils/chess/Types'

vi.mock('../../../src/utils/chess/Piece', () => {
    const Pawn = vi.fn()
    return { Pawn }
})

vi.mock('../../../src/utils/chess/Board', () => {
    const Board = vi.fn()
    return { Board }
})

describe('SimpleMove', () => {
    it('should get and set the correct pieces', () => {
        const mockedSetPiece = vi.fn()
        const mockedSetEnPassant = vi.fn()
        const mockedIncrement = vi.fn()
        const mockedDecrement = vi.fn()
        const board = new Board()
        const piece = new Pawn('', 0, 0)
        const newPiece = new Pawn('', 0, 0)
        const capturedPiece = new Pawn('', 0, 0)
        const enPassant = {x: 1, y: 1, xPiece: 2, yPiece: 2}
        board.setPiece = mockedSetPiece
        board.setEnPassant = mockedSetEnPassant
        board.increment = mockedIncrement
        board.decrement = mockedDecrement
        board.getPiece = vi.fn((x, y) => {
            if (x == 0 && y == 0) {return piece}
            if (x == 1 && y == 1) {return capturedPiece}
            return null
        })
        board.getEnPassant = vi.fn(() => enPassant)
        piece.copy = vi.fn(() => newPiece)
        Object.defineProperty(piece, 'color', {value: 'piece'})

        const boardMove: BoardMove = {xFrom: 0, yFrom: 0, xTo: 1, yTo: 1, options: {}}
        const move = new SimpleMove(board, boardMove)

        expect(move.execute()).toBe(true)
        expect(mockedSetPiece).toBeCalledWith(0, 0, null)
        expect(mockedSetPiece).toBeCalledWith(1, 1, newPiece)
        expect(mockedIncrement).toBeCalled()
        expect(mockedSetEnPassant).toBeCalledWith('piece', null)
        expect(move.undo()).toBe(true)
        expect(mockedSetPiece).toBeCalledWith(0, 0, piece)
        expect(mockedSetPiece).toBeCalledWith(1, 1, capturedPiece)
        expect(mockedDecrement).toBeCalled()
        expect(mockedSetEnPassant).toBeCalledWith('piece', enPassant)
    })
})

describe('RevealEnPassantMove', () => {
    it('should get and set the correct pieces', () => {
        const mockedSetEnPassant = vi.fn()
        const board = new Board()
        const piece = new Pawn('', 0, 0)
        const enPassant = {x: 2, y: 2, xPiece: 1, yPiece: 1}
        board.getPiece = vi.fn(() => piece)
        board.setPiece = vi.fn()
        board.getEnPassant = vi.fn()
        board.setEnPassant = mockedSetEnPassant
        board.increment = vi.fn()
        board.decrement = vi.fn()
        piece.copy = vi.fn(() => piece)
        Object.defineProperty(piece, 'color', {value: 'piece'})

        const boardMove: BoardMove = {xFrom: 0, yFrom: 0, xTo: 1, yTo: 1, options: {revealEnPassant: {x: 2, y: 2}}}
        const move = new RevealEnPassantMove(board, boardMove)

        expect(move.execute()).toBe(true)
        expect(mockedSetEnPassant.mock.calls.length).toBe(2)
        expect(mockedSetEnPassant.mock.calls.at(-1)).toEqual(['piece', enPassant])
        expect(move.undo()).toBe(true)
    })
})

describe('CaptureEnPassantMove', () => {
    it('should get and set the correct pieces', () => {
        const mockedSetPiece = vi.fn()
        const board = new Board()
        const pieceWhite = new Pawn('', 0, 0)
        const pieceBlack = new Pawn('', 0, 0)
        const pieceGrey = new Pawn('', 0, 0)
        const enPassantGrey = {x: 2, y: 2, xPiece: 1, yPiece: 1}
        const enPassantBlack = {x: 2, y: 2, xPiece: 3, yPiece: 3}
        pieceWhite.copy = vi.fn(() => pieceWhite)
        pieceBlack.copy = vi.fn(() => pieceBlack)
        pieceGrey.copy = vi.fn(() => pieceGrey)
        board.getPiece = vi.fn((x, y) => {
            if (x == 1 && y == 1) {return pieceGrey}
            if (x == 3 && y == 3) {return pieceBlack}
            return pieceWhite
        })
        board.setPiece = mockedSetPiece
        board.getEnPassant = vi.fn()
        board.setEnPassant = vi.fn()
        board.increment = vi.fn()
        board.decrement = vi.fn()
        board.getEnPassants = vi.fn((color, x, y) => {
            if (color === 'white' && x === 1 && y === 1) {return [enPassantGrey, enPassantBlack]}
            return []
        })
        Object.defineProperty(pieceWhite, 'color', {value: 'white'})

        const boardMove: BoardMove = {xFrom: 0, yFrom: 0, xTo: 1, yTo: 1, options: {captureEnPassant: true}}
        const move = new CaptureEnPassantMove(board, boardMove)

        expect(move.execute()).toBe(true)
        expect(mockedSetPiece).toBeCalledWith(1, 1, null)
        expect(mockedSetPiece).toBeCalledWith(3, 3, null)
        expect(move.undo()).toBe(true)
        expect(mockedSetPiece).toBeCalledWith(1, 1, pieceGrey)
        expect(mockedSetPiece).toBeCalledWith(3, 3, pieceBlack)
    })
})