import React, { useState, useRef } from 'react'
import ChessSquare from './ChessSquare'
import ChessActionBar from './ChessActionBar'
import { Board } from '../utils/chess/Board'
import { SimpleInvoker } from '../utils/chess/Invoker'
import { MoveFactory } from '../utils/chess/Move'

const boardStyle: React.CSSProperties = {
    display: 'flex',
    flexWrap: 'wrap',
    width: '500px',
    height: '500px',
}

const ChessBoard = (): JSX.Element => {
    const invoker = useRef(new SimpleInvoker())
    const board = useRef(new Board())
    const [squares, setSquares] = useState<any[][]>(board.current.squares)
    const [selected, setSelected] = useState<{x: number, y: number}>({x: -1, y: -1})
    const [destinations, setDestinations] = useState<{x: number, y: number, options: object}[]>([])

    const clearSelect = (): void => {
        setDestinations([])
        setSelected({x: -1, y: -1})
    }

    const executeMove = (xFrom: number, yFrom: number, xTo: number, yTo: number, options: object): void => {
        const move = MoveFactory.createMove(board.current, xFrom, yFrom, xTo, yTo, options)
        invoker.current.execute(move)
        setSquares(s => [...s])
        clearSelect()
    }

    const undoMove = (): void => {
        invoker.current.undo()
        setSquares(s => [...s])
        clearSelect()
    }

    const redoMove = (): void => {
        invoker.current.redo()
        setSquares(s => [...s])
        clearSelect()
    }

    const clickSquare = (x: number, y: number): void => {
        const sameSquare = selected.x === x && selected.y === y
        if (sameSquare) {
            clearSelect()
            return
        }
        const destination = destinations.find((d) => d.x === x && d.y === y)
        if (destination) {
            executeMove(selected.x, selected.y, x, y, destination.options)
            return
        }
        const moves = board.current.getMoves(x, y)
        setDestinations(moves)
        setSelected({x, y})
    }
 
    return (
        <div>
            <div style={boardStyle}>
                {squares.map((r) => r.map((s) =>
                    <ChessSquare
                        key={s.id}
                        square={s}
                        isSelected={selected.x === s.x && selected.y === s.y}
                        isDestination={destinations.some((d) => d.x === s.x && d.y === s.y)}
                        clickSquare={clickSquare}
                    />
                ))}
            </div>
            <ChessActionBar undo={undoMove} redo={redoMove} />
        </div>
    )
}

export default ChessBoard