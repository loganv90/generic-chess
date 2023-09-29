import { describe, expect, it, vi } from 'vitest'
import { Square } from '../../../src/utils/chess/Square'
import { Pawn } from '../../../src/utils/chess/Piece'

vi.mock('../../../src/utils/chess/Piece', () => {
    const Pawn = vi.fn()
    return { Pawn }
})

describe('Square', () => {
    it('should have null piece when not given character', () => {
        const square = new Square('0-0', 0, 0)
        expect(square.getPiece()).toBe(null)
    })

    it('should have pawn piece when given p character', () => {
        const square = new Square('0-0', 0, 0, 'p')
        expect(square.getPiece()).instanceOf(Pawn)
    })
})