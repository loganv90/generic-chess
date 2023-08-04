import { describe, it, expect } from "vitest"
import Board from "../../src/utils/ChessUtils"

describe('Board', () => {
    it('should have default piece arrangement with no arguments', () => {
        const board = new Board()
        expect(board.player).toBe(0)
        expect(board.currentPlayer).toBe(0)
        expect(board.players).toStrictEqual(['white', 'black'])
        expect(board.halfMoveClock).toBe(0)
        expect(board.fullMoveNumber).toBe(1)
        for (let i = 0; i < 8; i++) {
            expect(board.squares[1][i].piece.name).toBe('p')
            expect(board.squares[1][i].piece.color).toBe('black')
            expect(board.squares[0][i].piece.color).toBe('black')
            expect(board.squares[6][i].piece.name).toBe('p')
            expect(board.squares[6][i].piece.color).toBe('white')
            expect(board.squares[7][i].piece.color).toBe('white')
            expect(board.squares[2][i].piece).toBe(null)
            expect(board.squares[3][i].piece).toBe(null)
            expect(board.squares[4][i].piece).toBe(null)
            expect(board.squares[5][i].piece).toBe(null)
        }
        for (let i = 0; i < 8; i+=7) {
            expect(board.squares[i][0].piece.name).toBe('r')
            expect(board.squares[i][7].piece.name).toBe('r')
            expect(board.squares[i][1].piece.name).toBe('n')
            expect(board.squares[i][6].piece.name).toBe('n')
            expect(board.squares[i][2].piece.name).toBe('b')
            expect(board.squares[i][5].piece.name).toBe('b')
            expect(board.squares[i][3].piece.name).toBe('q')
            expect(board.squares[i][4].piece.name).toBe('k')
        }
    })

    it('should throw error with invalid FEN string', () => {
        expect(() => {
            new Board('', '')
        }).toThrow('Invalid FEN string')
    })

    it('getPieceMoves() should return correct moves for unmoved pawn', () => {
        const board = new Board()

        const eMoves = [{x: 4, y: 5}, {x: 4, y: 4}]
        const aMoves = board.getPieceMoves(4, 6)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for moved pawn', () => {
        const board = new Board()
        board.squares[6][4].piece.moved = true
        
        const eMoves = [{x: 4, y: 5}]
        const aMoves = board.getPieceMoves(4, 6)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for pawn en passant', () => {
        const board = new Board('white', 'rnbqkbnr/pppppppp/8/3pP3/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
        board.squares[3][4].piece.moved = true
        board.squares[3][3].piece.moved = true
        board.squares[3][3].piece.enPassant = {x: 3, y: 2}

        const eMoves = [{x:4, y:2}, {x:3, y:2}]
        const aMoves = board.getPieceMoves(4, 3)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for pawn capture', () => {
        const board = new Board('white', 'rnbqkbnr/pppppppp/3p4/4P3/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
        board.squares[3][4].piece.moved = true
        board.squares[2][3].piece.moved = true
        
        const eMoves = [{x:4, y:2}, {x:3, y:2}]
        const aMoves = board.getPieceMoves(4, 3)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for unmoved knight', () => {
        const board = new Board()

        const eMoves = [{x: 0, y: 5}, {x: 2, y: 5}]
        const aMoves = board.getPieceMoves(1, 7)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for moved knight', () => {
        const board = new Board('white', 'rnbqkbnr/pppppppp/8/4N3/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')

        const eMoves = [
            {x: 3, y: 1}, {x: 5, y: 1}, {x: 2, y: 2}, {x: 6, y: 2},
            {x: 2, y: 4}, {x: 6, y: 4}, {x: 3, y: 5}, {x: 5, y: 5}
        ]
        const aMoves = board.getPieceMoves(4, 3)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for unmoved bishop', () => {
        const board = new Board()

        const eMoves = []
        const aMoves = board.getPieceMoves(2, 7)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for moved bishop', () => {
        const board = new Board('white', 'rnbqkbnr/pppppppp/8/6B1/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')

        const eMoves = [
            {x: 7, y: 4}, {x: 7, y: 2}, {x: 5, y: 2}, {x: 4, y: 1},
            {x: 5, y: 4}, {x: 4, y: 5}
        ]
        const aMoves = board.getPieceMoves(6, 3)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for unmoved rook', () => {
        const board = new Board()

        const eMoves = []
        const aMoves = board.getPieceMoves(0, 7)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for moved rook', () => {
        const board = new Board('white', 'rnbqkbnr/pppppppp/8/8/R7/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')

        const eMoves = [
            {x: 0, y: 5}, {x: 0, y: 3}, {x: 0, y: 2}, {x: 0, y: 1},
            {x: 1, y: 4}, {x: 2, y: 4}, {x: 3, y: 4}, {x: 4, y: 4},
            {x: 5, y: 4}, {x: 6, y: 4}, {x: 7, y: 4}
        ]
        const aMoves = board.getPieceMoves(0, 4)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for unmoved queen', () => {
        const board = new Board()

        const eMoves = []
        const aMoves = board.getPieceMoves(3, 7)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for moved queen', () => {
        const board = new Board('white', 'rnbqkbnr/pppppppp/8/8/3Q4/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')

        const eMoves = [
            {x: 4, y: 5}, {x: 4, y: 3}, {x: 4, y: 4}, {x: 5, y: 4},
            {x: 5, y: 2}, {x: 6, y: 1}, {x: 6, y: 4}, {x: 7, y: 4},
            {x: 2, y: 5}, {x: 2, y: 3}, {x: 2, y: 4}, {x: 1, y: 4},
            {x: 1, y: 2}, {x: 0, y: 1}, {x: 0, y: 4}, {x: 3, y: 5},
            {x: 3, y: 3}, {x: 3, y: 2}, {x: 3, y: 1}
        ]
        const aMoves = board.getPieceMoves(3, 4)

        console.log(aMoves)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for unmoved king', () => {
        const board = new Board()

        const eMoves = []
        const aMoves = board.getPieceMoves(4, 7)

        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })

    it('getPieceMoves() should return correct moves for moved king', () => {
        const board = new Board('white', 'rnbqkbnr/pppppppp/8/8/4K3/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')

        const eMoves = [
            {x: 5, y: 5}, {x: 5, y: 4},
            {x: 5, y: 3}, {x: 3, y: 4},
            {x: 3, y: 5}, {x: 4, y: 5},
            {x: 3, y: 3}, {x: 4, y: 3}
        ]
        const aMoves = board.getPieceMoves(4, 4)

        console.log(aMoves)


        expect(aMoves).toEqual(expect.arrayContaining(eMoves))
        expect(aMoves.length).toBe(eMoves.length)
    })
})
