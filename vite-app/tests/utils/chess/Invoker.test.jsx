import { describe, it, expect, vi } from 'vitest'
import { SimpleInvoker } from '/src/utils/chess/Invoker'
import { Move } from '/src/utils/chess/Move'

vi.mock('/src/utils/chess/Move', () => {
    const Move = vi.fn(() => ({
        execute: vi.fn(() => true),
        undo: vi.fn(() => true),
    }))
    return { Move }
})

describe('SimpleInvoker', () => {
    it('should undo and redo moves in order', () => {
        const simpleInvoker = new SimpleInvoker()
        const move1 = new Move()
        const move2 = new Move()

        expect(simpleInvoker.execute(move1)).toBe(true)
        expect(move1.execute).toBeCalledTimes(1)

        expect(simpleInvoker.execute(move2)).toBe(true)
        expect(move2.execute).toBeCalledTimes(1)

        expect(simpleInvoker.undo()).toBe(true)
        expect(move2.undo).toBeCalledTimes(1)

        expect(simpleInvoker.undo()).toBe(true)
        expect(move1.undo).toBeCalledTimes(1)

        expect(simpleInvoker.redo()).toBe(true)
        expect(move1.execute).toBeCalledTimes(2)

        expect(simpleInvoker.redo()).toBe(true)
        expect(move2.execute).toBeCalledTimes(2)
    })

    it('should not undo or redo moves if there are no moves', () => {
        const simpleInvoker = new SimpleInvoker()

        expect(simpleInvoker.undo()).toBe(false)
        expect(simpleInvoker.redo()).toBe(false)
    })

    it('should overwrite history when executing a new move', () => {
        const simpleInvoker = new SimpleInvoker()

        const move1 = new Move()
        const move2 = new Move()
        const move3 = new Move()

        expect(simpleInvoker.execute(move1)).toBe(true)
        expect(simpleInvoker.execute(move2)).toBe(true)

        expect(simpleInvoker.undo()).toBe(true)
        expect(simpleInvoker.undo()).toBe(true)

        expect(simpleInvoker.execute(move3)).toBe(true)
        expect(simpleInvoker.redo()).toBe(false)
    })
})
