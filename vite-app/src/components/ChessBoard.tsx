import React, { useState, useRef } from 'react'
import ChessSquare from './ChessSquare'
import ChessActionBar from './ChessActionBar'
import { Board } from '../utils/chess/Board'
import { SimpleInvoker } from '../utils/chess/Invoker'
import { Square } from '../utils/chess/Square'
import { Move } from '../utils/chess/Move'

const boardStyle: React.CSSProperties = {
    display: 'flex',
    flexWrap: 'wrap',
    width: '500px',
    height: '500px',
}

const ChessBoard = (): JSX.Element => {
    const invoker = useRef(new SimpleInvoker())
    const board = useRef(new Board())
    const [squares, setSquares] = useState<Square[][]>(board.current.squares)
    const [selected, setSelected] = useState<{x: number, y: number}>({x: -1, y: -1})
    const [moves, setMoves] = useState<Move[]>([])

    const clearSelect = (): void => {
        setMoves([])
        setSelected({x: -1, y: -1})
    }

    const executeMove = (move: Move): void => {
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

        const move = moves.find(m => m.xTo === x && m.yTo === y)
        if (move) {
            executeMove(move)
            return
        }

        setMoves(board.current.getMoves(x, y))
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
                        isDestination={moves.some(m => m.xTo === s.x && m.yTo === s.y)}
                        clickSquare={clickSquare}
                    />
                ))}
            </div>
            <ChessActionBar undo={undoMove} redo={redoMove} />
        </div>
    )
}

export default ChessBoard