import { useState, useRef } from 'react'
import ChessSquare from '/src/components/ChessSquare'
import ChessActionBar from '/src/components/ChessActionBar'
import { Board } from '/src/utils/chess/Board'
import { SimpleInvoker } from '/src/utils/chess/Invoker'
import { MoveFactory } from '/src/utils/chess/Move'

function ChessBoard() {
    const invoker = useRef(new SimpleInvoker())
    const board = useRef(new Board())
    const [squares, setSquares] = useState(board.current.squares)
    const [selected, setSelected] = useState({})
    const [destinations, setDestinations] = useState([])

    const clearSelect = () => {
        setDestinations([])
        setSelected({})
    }

    const executeMove = (xFrom, yFrom, xTo, yTo, options) => {
        const move = MoveFactory.createMove(board.current, xFrom, yFrom, xTo, yTo, options)
        invoker.current.execute(move)
        setSquares(s => [...s])
        clearSelect()
    }

    const undoMove = () => {
        invoker.current.undo()
        setSquares(s => [...s])
        clearSelect()
    }

    const redoMove = () => {
        invoker.current.redo()
        setSquares(s => [...s])
        clearSelect()
    }

    const clickSquare = (x, y) => {
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
            <div style={{...boardStyle}}>
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

const boardStyle = {
    display: 'flex',
    flexWrap: 'wrap',
    width: '500px',
    height: '500px',
}

export default ChessBoard
