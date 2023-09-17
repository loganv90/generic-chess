import { Move } from './Move'

abstract class Invoker {
    abstract execute(move: Move): boolean
    abstract undo(): boolean
    abstract redo(): boolean
}

class SimpleInvoker extends Invoker {
    private history: Move[]
    private index: number

    constructor() {
        super()
        this.history = []
        this.index = 0
    }

    execute(move: Move): boolean {
        const success = move.execute()

        if (success) {
            this.history.splice(this.index, this.history.length-this.index, move)
            this.index += 1
        } else {
            throw new Error('Failed to execute move')
        }

        return success
    }

    undo() {
        if (this.index <= 0) {
            return false
        }

        const move = this.history[this.index-1]
        const success = move.undo()

        if (success) {
            this.index -= 1
        } else {
            throw new Error('Failed to undo move')
        }

        return success
    }

    redo() {
        if (this.index >= this.history.length) {
            return false
        }

        const move = this.history[this.index]
        const success = move.execute()

        if (success) {
            this.index += 1
        } else {
            throw new Error('Failed to redo move')
        }

        return success
    }
}

export { Invoker, SimpleInvoker }
