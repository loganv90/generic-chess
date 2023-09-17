import React, { useState, useRef } from 'react'
import ChessSquare from './ChessSquare'
import ChessActionBar from './ChessActionBar'
import { Board } from '../utils/chess/Board'
import { SimpleInvoker } from '../utils/chess/Invoker'
import { MoveFactory } from '../utils/chess/Move'
import { Square } from '../utils/chess/Square'
import { BoardMove } from '../utils/chess/Types'
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
    const [boardMoves, setBoardMoves] = useState<BoardMove[]>([])

    const clearSelect = (): void => {
        setBoardMoves([])
        setSelected({x: -1, y: -1})
    }

    const executeMove = (boardMove: BoardMove): void => {
        const move: Move | null = MoveFactory.createMove(board.current, boardMove)
        move && invoker.current.execute(move)
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

        const boardMove = boardMoves.find((d) => d.xTo === x && d.yTo === y)
        if (boardMove) {
            executeMove(boardMove)
            return
        }

        const moves = board.current.getMoves(x, y)
        setBoardMoves(moves)
        setSelected({x, y})

        console.log(boardMoves)
    }
 
    return (
        <div>
            <div style={boardStyle}>
                {squares.map((r) => r.map((s) =>
                    <ChessSquare
                        key={s.id}
                        square={s}
                        isSelected={selected.x === s.x && selected.y === s.y}
                        isDestination={boardMoves.some((d) => d.xTo === s.x && d.yTo === s.y)}
                        clickSquare={clickSquare}
                    />
                ))}
            </div>
            <ChessActionBar undo={undoMove} redo={redoMove} />
        </div>
    )
}

export default ChessBoard