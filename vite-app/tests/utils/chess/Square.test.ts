import { describe, expect, it, vi } from 'vitest'
import { Square } from '../../../src/utils/chess/Square'
import { PieceMove } from '../../../src/utils/chess/Types'

vi.mock('../../../src/utils/chess/Piece', () => {
    const Pawn = vi.fn(() => ({
        getMoves: vi.fn((): PieceMove[] => [{x: 0, y: 0, options: {}}])
    }))
    return { Pawn }
})

describe('Square', () => {
    describe('getMoves', () => {
        it('should return empty array if there is no piece', () => {
            const square = new Square('0-0', 0, 0)
            expect(square.getMoves().length).toBe(0)
        })

        it('should return array if there is a piece', () => {
            const square = new Square('0-0', 0, 0, 'p')
            expect(square.getMoves().length).toBeGreaterThan(0)
        })
    })
})