import { useState } from 'react'
import PropTypes from 'prop-types'
import ChessSquare from './ChessSquare'

function ChessBoard({ xSquares, ySquares }) {
    const createSquares = (numX, numY) => {
        const s = []
        const squareWidth = `${100 / numX}%`
        const squareHeight = `${100 / numY}%`
        for (let y = 0; y < numY; y++) {
            for (let x = 0; x < numX; x++) {
                let id = `${x}${y}`
                let isDark = (x + y) % 2 === 0

                s.push({
                    id: id,
                    x: x,
                    y: y,
                    isDark: isDark,
                    width: squareWidth,
                    height: squareHeight,
                    isSelected: false,
                })
            }
        }
        return s
    }

    const [squares, setSquares] = useState(createSquares(xSquares, ySquares))
    const [selectedSquare, setSelectedSquare] = useState('')

    const clickSquare = (id) => {
        let newSelectedSquare = selectedSquare

        if (newSelectedSquare === id) {
            newSelectedSquare = ''
        } else if (newSelectedSquare === '') {
            newSelectedSquare = id
        } else {
            console.log(`Move from ${selectedSquare} to ${id}`)
            newSelectedSquare = ''
        }

        const newSquares = squares.map((s) => {
            return {...s, isSelected: s.id === newSelectedSquare}
        })

        setSelectedSquare(newSelectedSquare)
        setSquares(newSquares)
    }
 
    return (
        <div style={{...boardStyle}}>
            {squares.map((s) => (
                <ChessSquare key={s.id} square={s} clickSquare={clickSquare} />
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

ChessBoard.defaultProps = {
    xSquares: 8,
    ySquares: 8,
}

ChessBoard.propTypes = {
    xSquares: PropTypes.number,
    ySquares: PropTypes.number,
}

export default ChessBoard
