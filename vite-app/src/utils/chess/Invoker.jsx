class Invoker {
    constructor() {}

    execute() {
        throw new Error('Invoker.execute not implemented')
    }

    undo() {
        throw new Error('Invoker.undo not implemented')
    }

    redo() {
        throw new Error('Invoker.redo not implemented')
    }
}

class SimpleInvoker extends Invoker {
    constructor() {
        super()
        this.history = []
        this.index = 0
    }

    execute(move) {
        const success = move.execute()

        if (success) {
            this.history.splice(this.index, this.history.length-this.index, move)
            this.index += 1
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
        }

        return success
    }
}

export { Invoker, SimpleInvoker }
