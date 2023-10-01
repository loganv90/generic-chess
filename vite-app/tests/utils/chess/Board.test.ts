import { describe, it, expect } from 'vitest'
import { Board } from '../../../src/utils/chess/Board'
import { Pawn, Knight, Bishop, Rook, Queen, King } from '../../../src/utils/chess/Piece'

describe('Board', () => {
    it('should have default pieces', () => {
        const board = new Board()
        for (let i=0; i<8; i+=7) {
            expect(board.getPiece(0, i) instanceof Rook).toBe(true)
            expect(board.getPiece(1, i) instanceof Knight).toBe(true)
            expect(board.getPiece(2, i) instanceof Bishop).toBe(true)
            expect(board.getPiece(3, i) instanceof Queen).toBe(true)
            expect(board.getPiece(4, i) instanceof King).toBe(true)
            expect(board.getPiece(5, i) instanceof Bishop).toBe(true)
            expect(board.getPiece(6, i) instanceof Knight).toBe(true)
            expect(board.getPiece(7, i) instanceof Rook).toBe(true)
        }
        for (let i=0; i<8; i++) {
            expect(board.getPiece(i, 1) instanceof Pawn).toBe(true)
            expect(board.getPiece(i, 6) instanceof Pawn).toBe(true)
            expect(board.getPiece(i, 2)).toBe(null)
            expect(board.getPiece(i, 3)).toBe(null)
            expect(board.getPiece(i, 4)).toBe(null)
            expect(board.getPiece(i, 5)).toBe(null)
            for (let j=0; j<2; j++) {
                expect(board.getPiece(0, j)?.color).toBe('black')
            }
            for (let j=6; j<8; j++) {
                expect(board.getPiece(0, j)?.color).toBe('white')
            }
        }
    })

    // describe('getMoves', () => {
    //     it('should return correct moves for unmoved pawn', () => {
    //         const board = new Board()
    
    //         const eMoves = [
    //             {xFrom: 4, yFrom: 6, xTo: 4, yTo: 5, options: {}},
    //             {xFrom: 4, yFrom: 6, xTo: 4, yTo: 4, options: {revealEnPassant: {x: 4, y: 5}}},
    //         ]
    //         const aMoves = board.getMoves(4, 6)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for moved pawn', () => {
    //         const board = new Board()
    //         board.setPiece(4, 6, new Pawn('white', 0, -1, true))
            
    //         const eMoves = [
    //             {xFrom: 4, yFrom: 6, xTo: 4, yTo: 5, options: {}},
    //         ]
    //         const aMoves = board.getMoves(4, 6)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for pawn en passant', () => {
    //         const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/8/3pP3/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    //         board.setPiece(4, 3, new Pawn('white', 0, -1, true))
    //         board.setEnPassant('black', {x: 3, y: 2, xPiece: 3, yPiece: 3})
    
    //         const eMoves = [
    //             {xFrom: 4, yFrom: 3, xTo: 4, yTo: 2, options: {}},
    //             {xFrom: 4, yFrom: 3, xTo: 3, yTo: 2, options: {captureEnPassant: true}},
    //         ]
    //         const aMoves = board.getMoves(4, 3)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for pawn capture', () => {
    //         const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/3p4/4P3/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    //         board.setPiece(4, 3, new Pawn('white', 0, -1, true))
            
    //         const eMoves = [
    //             {xFrom: 4, yFrom: 3, xTo: 4, yTo: 2, options: {}},
    //             {xFrom: 4, yFrom: 3, xTo: 3, yTo: 2, options: {}},
    //         ]
    //         const aMoves = board.getMoves(4, 3)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })

    //     it('should return correct moves for pawn promotion', () => {
    //         const board = new Board(0, 0, ['white', 'black'], 'rnbqkbn1/pppppppP/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    //         board.setPiece(7, 1, new Pawn('white', 0, -1, true))

    //         const eMoves = [
    //             {xFrom: 7, yFrom: 1, xTo: 7, yTo: 0, options: {promote: true}},
    //             {xFrom: 7, yFrom: 1, xTo: 6, yTo: 0, options: {promote: true}},
    //         ]
    //         const aMoves = board.getMoves(7, 1)

    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for unmoved knight', () => {
    //         const board = new Board()
    
    //         const eMoves = [
    //             {xFrom: 1, yFrom: 7, xTo: 0, yTo: 5, options: {}},
    //             {xFrom: 1, yFrom: 7, xTo: 2, yTo: 5, options: {}},
    //         ]
    //         const aMoves = board.getMoves(1, 7)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for moved knight', () => {
    //         const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/8/4N3/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    
    //         const eMoves = [
    //             {xFrom: 4, yFrom: 3, xTo: 3, yTo: 1, options: {}},
    //             {xFrom: 4, yFrom: 3, xTo: 5, yTo: 1, options: {}},
    //             {xFrom: 4, yFrom: 3, xTo: 2, yTo: 2, options: {}},
    //             {xFrom: 4, yFrom: 3, xTo: 6, yTo: 2, options: {}},
    //             {xFrom: 4, yFrom: 3, xTo: 2, yTo: 4, options: {}},
    //             {xFrom: 4, yFrom: 3, xTo: 6, yTo: 4, options: {}},
    //             {xFrom: 4, yFrom: 3, xTo: 3, yTo: 5, options: {}},
    //             {xFrom: 4, yFrom: 3, xTo: 5, yTo: 5, options: {}},
    //         ]
    //         const aMoves = board.getMoves(4, 3)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for unmoved bishop', () => {
    //         const board = new Board()
    
    //         const eMoves: PieceMove[] = []
    //         const aMoves = board.getMoves(2, 7)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for moved bishop', () => {
    //         const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/8/6B1/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    
    //         const eMoves = [
    //             {xFrom: 6, yFrom: 3, xTo: 7, yTo: 4, options: {}},
    //             {xFrom: 6, yFrom: 3, xTo: 7, yTo: 2, options: {}},
    //             {xFrom: 6, yFrom: 3, xTo: 5, yTo: 2, options: {}},
    //             {xFrom: 6, yFrom: 3, xTo: 4, yTo: 1, options: {}},
    //             {xFrom: 6, yFrom: 3, xTo: 5, yTo: 4, options: {}},
    //             {xFrom: 6, yFrom: 3, xTo: 4, yTo: 5, options: {}},
    //         ]
    //         const aMoves = board.getMoves(6, 3)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for unmoved rook', () => {
    //         const board = new Board()
    
    //         const eMoves: PieceMove[] = []
    //         const aMoves = board.getMoves(0, 7)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for moved rook', () => {
    //         const board = new Board(0, 0, ['white','black'], 'rnbqkbnr/pppppppp/8/8/R7/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    
    //         const eMoves = [
    //             {xFrom: 0, yFrom: 4, xTo: 0, yTo: 5, options: {}},
    //             {xFrom: 0, yFrom: 4, xTo: 0, yTo: 3, options: {}},
    //             {xFrom: 0, yFrom: 4, xTo: 0, yTo: 2, options: {}},
    //             {xFrom: 0, yFrom: 4, xTo: 0, yTo: 1, options: {}},
    //             {xFrom: 0, yFrom: 4, xTo: 1, yTo: 4, options: {}},
    //             {xFrom: 0, yFrom: 4, xTo: 2, yTo: 4, options: {}},
    //             {xFrom: 0, yFrom: 4, xTo: 3, yTo: 4, options: {}},
    //             {xFrom: 0, yFrom: 4, xTo: 4, yTo: 4, options: {}},
    //             {xFrom: 0, yFrom: 4, xTo: 5, yTo: 4, options: {}},
    //             {xFrom: 0, yFrom: 4, xTo: 6, yTo: 4, options: {}},
    //             {xFrom: 0, yFrom: 4, xTo: 7, yTo: 4, options: {}},
    //         ]
    //         const aMoves = board.getMoves(0, 4)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for unmoved queen', () => {
    //         const board = new Board()
    
    //         const eMoves: PieceMove[] = []
    //         const aMoves = board.getMoves(3, 7)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for moved queen', () => {
    //         const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/8/8/3Q4/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    
    //         const eMoves = [
    //             {xFrom: 3, yFrom: 4, xTo: 4, yTo: 5, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 4, yTo: 3, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 4, yTo: 4, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 5, yTo: 4, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 5, yTo: 2, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 6, yTo: 1, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 6, yTo: 4, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 7, yTo: 4, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 2, yTo: 5, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 2, yTo: 3, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 2, yTo: 4, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 1, yTo: 4, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 1, yTo: 2, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 0, yTo: 1, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 0, yTo: 4, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 3, yTo: 5, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 3, yTo: 3, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 3, yTo: 2, options: {}},
    //             {xFrom: 3, yFrom: 4, xTo: 3, yTo: 1, options: {}},
    //         ]
    //         const aMoves = board.getMoves(3, 4)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for unmoved king', () => {
    //         const board = new Board()
    
    //         const eMoves: PieceMove[] = []
    //         const aMoves = board.getMoves(4, 7)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    
    //     it('should return correct moves for moved king', () => {
    //         const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/8/8/4K3/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    
    //         const eMoves = [
    //             {xFrom: 4, yFrom: 4, xTo: 5, yTo: 5, options: {}},
    //             {xFrom: 4, yFrom: 4, xTo: 5, yTo: 4, options: {}},
    //             {xFrom: 4, yFrom: 4, xTo: 5, yTo: 3, options: {}},
    //             {xFrom: 4, yFrom: 4, xTo: 3, yTo: 4, options: {}},
    //             {xFrom: 4, yFrom: 4, xTo: 3, yTo: 5, options: {}},
    //             {xFrom: 4, yFrom: 4, xTo: 4, yTo: 5, options: {}},
    //             {xFrom: 4, yFrom: 4, xTo: 3, yTo: 3, options: {}},
    //             {xFrom: 4, yFrom: 4, xTo: 4, yTo: 3, options: {}},
    //             {xFrom: 4, yFrom: 4, xTo: 5, yTo: 4, options: {}}, // castle one direction
    //             {xFrom: 4, yFrom: 4, xTo: 3, yTo: 4, options: {}}, // castle other direction
    //             // also make a test for moved king and rooks
    //         ]
    //         const aMoves = board.getMoves(4, 4)
    
    //         expect(aMoves).toEqual(expect.arrayContaining(eMoves))
    //         expect(aMoves.length).toBe(eMoves.length)
    //     })
    // })

    // describe('getEnPassants', () => {
    //     it('should return correct en passants', () => {
    //         const board = new Board()
            
    //         const enPassantWhite = {x: 1, y: 1, xPiece: 2, yPiece: 2}
    //         const enPassantBlack = {x: 1, y: 1, xPiece: 3, yPiece: 3}
    //         const enPassantGrey = {x: 1, y: 1, xPiece: 4, yPiece: 4}
    //         board.setEnPassant('white', enPassantWhite)
    //         board.setEnPassant('black', enPassantBlack)
    //         board.setEnPassant('grey', enPassantGrey)

    //         const eEnPassants = [enPassantBlack, enPassantGrey]
    //         const aEnPassants = board.getEnPassants('white', 1, 1)

    //         expect(aEnPassants).toEqual(expect.arrayContaining(eEnPassants))
    //         expect(aEnPassants.length).toBe(eEnPassants.length)
    //     })
    // })
})