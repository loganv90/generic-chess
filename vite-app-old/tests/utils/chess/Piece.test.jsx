import { describe, it } from 'vitest'
import { Pawn, Knight, Bishop, Rook, Queen, King } from '/src/utils/chess/Piece'

describe('Pawn', () => {
    it('should be able to be created', () => {
        new Pawn('white', 0, -1)
    })
})

describe('Knight', () => {
    it('should be able to be created', () => {
        new Knight('white')
    })
})

describe('Bishop', () => {
    it('should be able to be created', () => {
        new Bishop('white')
    })
})

describe('Rook', () => {
    it('should be able to be created', () => {
        new Rook('white')
    })
})

describe('Queen', () => {
    it('should be able to be created', () => {
        new Queen('white')
    })
})

describe('King', () => {
    it('should be able to be created', () => {
        new King('white')
    })
})
