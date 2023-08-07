import { useState } from 'react'
import ChessSquare from '/src/components/ChessSquare'
import { Board } from '/src/utils/chess/Board'

function ChessBoard() {
    const [board, setBoard] = useState(new Board())
    const [selected, setSelected] = useState({})
    const [destinations, setDestinations] = useState([])

    const clickSquare = (x, y) => {
        if (selected.x === x && selected.y === y) {
            setDestinations([])
            setSelected({})
        } else if (destinations.some((d) => d.x === x && d.y === y)) {
            board.movePiece(selected.x, selected.y, x, y)
            setBoard(board)
            setDestinations([])
            setSelected({})
        } else {
            setDestinations(board.getPieceMoves(x, y))
            setSelected({x, y})
        }
    }
 
    return (
        <div style={{...boardStyle}}>
            {board.squares.map((row) => row.map((s) =>
                <ChessSquare
                    key={s.id}
                    square={s}
                    isSelected={selected.x === s.x && selected.y === s.y}
                    isDestination={destinations.some((d) => d.x === s.x && d.y === s.y)}
                    clickSquare={clickSquare}
                />
            ))}
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
