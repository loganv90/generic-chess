import { describe, it } from 'vitest'
import { Square } from '/src/utils/chess/Square'

describe('Square', () => {
    it('should be able to be created', () => {
        new Square('0-0', 0, 0, true, 'r')
    })
})
