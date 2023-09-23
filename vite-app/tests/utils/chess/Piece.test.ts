import { describe, it, expect } from 'vitest'
import { Pawn, Knight, Bishop, Rook, Queen, King } from '../../../src/utils/chess/Piece'

describe('Pawn', () => {
    describe('getMoves', () => {
        it('should return pawn moves', () => {
            const pawn1 = new Pawn('white', 0, -1)

            const eMoves1 = [
                {x: 0, y: -1, options: {canPromote: true, noCapture: true}},
                {x: 0, y: -2, options: {canPromote: true, noCapture: true, mustRevealEnPassant: {x: 0, y: -1}}},
                {x: 1, y: -1, options: {canPromote: true, mustCapture: true, canCaptureEnPassant: true}},
                {x: -1, y: -1, options: {canPromote: true, mustCapture: true, canCaptureEnPassant: true}}
            ]
            const aMoves1 = pawn1.getMoves()

            expect(aMoves1).toEqual(expect.arrayContaining(eMoves1))
            expect(aMoves1.length).toBe(eMoves1.length)

            const pawn2 = pawn1.copy()

            const eMoves2 = [
                {x: 0, y: -1, options: {canPromote: true, noCapture: true}},
                {x: 1, y: -1, options: {canPromote: true, mustCapture: true, canCaptureEnPassant: true}},
                {x: -1, y: -1, options: {canPromote: true, mustCapture: true, canCaptureEnPassant: true}}
            ]
            const aMoves2 = pawn2.getMoves()

            expect(aMoves2).toEqual(expect.arrayContaining(eMoves2))
            expect(aMoves2.length).toBe(eMoves2.length)
        })
    })
})

describe('Knight', () => {
    describe('getMoves', () => {
        it('should return knight moves', () => {
            const knight = new Knight('white')

            const eMoves = [
                {x: 1, y: 2, options: {}},
                {x: 2, y: 1, options: {}},
                {x: 2, y: -1, options: {}},
                {x: 1, y: -2, options: {}},
                {x: -1, y: -2, options: {}},
                {x: -2, y: -1, options: {}},
                {x: -2, y: 1, options: {}},
                {x: -1, y: 2, options: {}},
            ]
            const aMoves = knight.getMoves()

            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    })
})

describe('Bishop', () => {
    describe('getMoves', () => {
        it('should return bishop moves', () => {
            const bishop = new Bishop('white')

            const eMoves = [
                {x: 1, y: 1, options: {direction: true}},
                {x: 1, y: -1, options: {direction: true}},
                {x: -1, y: 1, options: {direction: true}},
                {x: -1, y: -1, options: {direction: true}},
            ]
            const aMoves = bishop.getMoves()

            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    })
})

describe('Rook', () => {
    describe('getMoves', () => {
        it('should return rook moves', () => {
            const rook = new Rook('white')

            const eMoves = [
                {x: 1, y: 0, options: {direction: true}},
                {x: 0, y: 1, options: {direction: true}},
                {x: -1, y: 0, options: {direction: true}},
                {x: 0, y: -1, options: {direction: true}},
            ]
            const aMoves = rook.getMoves()

            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    })
})

describe('Queen', () => {
    describe('getMoves', () => {
        it('should return queen moves', () => {
            const queen = new Queen('white')

            const eMoves = [
                {x: 1, y: 1, options: {direction: true}},
                {x: 1, y: -1, options: {direction: true}},
                {x: -1, y: 1, options: {direction: true}},
                {x: -1, y: -1, options: {direction: true}},
                {x: 1, y: 0, options: {direction: true}},
                {x: 0, y: 1, options: {direction: true}},
                {x: -1, y: 0, options: {direction: true}},
                {x: 0, y: -1, options: {direction: true}},
            ]
            const aMoves = queen.getMoves()

            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    })
})

describe('King', () => {
    describe('getMoves', () => {
        it('should return king moves', () => {
            const king = new King('white', 0, 1)

            const eMoves = [
                {x: 1, y: 1, options: {}},
                {x: 1, y: 0, options: {}},
                {x: 1, y: -1, options: {}},
                {x: 0, y: 1, options: {}},
                {x: 0, y: -1, options: {}},
                {x: -1, y: 1, options: {}},
                {x: -1, y: 0, options: {}},
                {x: -1, y: -1, options: {}},
                {x: 1, y: 0, options: {mustCastle: true}},
                {x: -1, y: 0, options: {mustCastle: true}},
            ]
            const aMoves = king.getMoves()

            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    })
})