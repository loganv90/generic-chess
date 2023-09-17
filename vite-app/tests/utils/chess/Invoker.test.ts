import { describe, it, expect, vi } from 'vitest'
import { SimpleInvoker } from '../../../src/utils/chess/Invoker'
import { SimpleMove } from '../../../src/utils/chess/Move'
import { Board } from '../../../src/utils/chess/Board'
import { BoardMove } from '../../../src/utils/chess/Types'

vi.mock('/src/utils/chess/Move', () => {
    const SimpleMove = vi.fn(() => ({
        execute: vi.fn(() => true),
        undo: vi.fn(() => true),
    }))
    return { SimpleMove }
})

describe('SimpleInvoker', () => {
    it('should undo and redo moves in order', () => {
        const mockedExecute1 = vi.fn(() => true)
        const mockedUndo1 = vi.fn(() => true)
        const mockedExecute2 = vi.fn(() => true)
        const mockedUndo2 = vi.fn(() => true)
        const simpleInvoker = new SimpleInvoker()
        const board = new Board()
        const boardMove: BoardMove = {xFrom: 0, yFrom: 0, xTo: 0, yTo: 0, options: {}}
        const move1 = new SimpleMove(board, boardMove)
        const move2 = new SimpleMove(board, boardMove)
        move1.execute = mockedExecute1
        move1.undo = mockedUndo1
        move2.execute = mockedExecute2
        move2.undo = mockedUndo2

        expect(simpleInvoker.execute(move1)).toBe(true)
        expect(mockedExecute1).toBeCalledTimes(1)

        expect(simpleInvoker.execute(move2)).toBe(true)
        expect(mockedExecute2).toBeCalledTimes(1)

        expect(simpleInvoker.undo()).toBe(true)
        expect(mockedUndo2).toBeCalledTimes(1)

        expect(simpleInvoker.undo()).toBe(true)
        expect(mockedUndo1).toBeCalledTimes(1)

        expect(simpleInvoker.redo()).toBe(true)
        expect(mockedExecute1).toBeCalledTimes(2)

        expect(simpleInvoker.redo()).toBe(true)
        expect(mockedExecute2).toBeCalledTimes(2)
    })

    it('should not undo or redo moves if there are no moves', () => {
        const simpleInvoker = new SimpleInvoker()

        expect(simpleInvoker.undo()).toBe(false)
        expect(simpleInvoker.redo()).toBe(false)
    })

    it('should overwrite history when executing a new move', () => {
        const simpleInvoker = new SimpleInvoker()
        const board = new Board()
        const boardMove: BoardMove = {xFrom: 0, yFrom: 0, xTo: 0, yTo: 0, options: {}}
        const move1 = new SimpleMove(board, boardMove)
        const move2 = new SimpleMove(board, boardMove)
        const move3 = new SimpleMove(board, boardMove)

        expect(simpleInvoker.execute(move1)).toBe(true)
        expect(simpleInvoker.execute(move2)).toBe(true)

        expect(simpleInvoker.undo()).toBe(true)
        expect(simpleInvoker.undo()).toBe(true)

        expect(simpleInvoker.execute(move3)).toBe(true)
        expect(simpleInvoker.redo()).toBe(false)
    })
})
