import { describe, it, expect } from "vitest"
import getPieceMoves from "../../src/utils/ChessUtils"

describe('getPieceMoves', () => {
    it('should return an empty array', () => {
        expect(getPieceMoves()).toEqual([])
    })
})

export default getPieceMoves
