import { useState } from 'react'
import PropTypes from 'prop-types'
import ChessSquare from './ChessSquare'

function ChessBoard({ boardConfig }) {
    const createSquareConfigs = (numX, numY) => {
        const squareConfigs = []
        const squareWidth = `${100 / numX}%`
        const squareHeight = `${100 / numY}%`
        for (let y = 0; y < numY; y++) {
            for (let x = 0; x < numX; x++) {
                let id = `${x}${y}`
                let isDark = (x + y) % 2 === 0
                let piece = boardConfig.startingPieces[id] ? boardConfig.startingPieces[id] : null

                squareConfigs.push({
                    id: id,
                    x: x,
                    y: y,
                    isDark: isDark,
                    width: squareWidth,
                    height: squareHeight,
                    isSelected: false,
                    isDestination: false,
                    piece: piece,
                })
            }
        }
        return squareConfigs
    }

    const [squareConfigs, setSquareConfigs] = useState(createSquareConfigs(boardConfig.xSquares, boardConfig.ySquares))
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

        const newSquares = squareConfigs.map((s) => {
            return {...s, isSelected: s.id === newSelectedSquare}
        })

        setSelectedSquare(newSelectedSquare)
        setSquareConfigs(newSquares)
    }
 
    return (
        <div style={{...boardStyle}}>
            {squareConfigs.map((s) => (
                <ChessSquare key={s.id} squareConfig={s} clickSquare={clickSquare} />
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
    boardConfig: {},
}

ChessBoard.propTypes = {
    boardConfig: PropTypes.object,
}

export default ChessBoard
