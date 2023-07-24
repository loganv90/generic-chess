import { useState } from 'react'
import PropTypes from 'prop-types'

function ChessSquare({ square, clickSquare }) {
    const getSquareColor = () => {
        if (square.isSelected) {
            return 'red'
        } else if (square.isDark) {
            return 'darkGrey'
        } else {
            return 'lightGrey'
        }
    }

    return (
        <div
            style={{
                ...squareStyle,
                backgroundColor: getSquareColor(),
                width: square.width,
                height: square.height
            }}
            onClick={() => clickSquare(square.id)}
        >
            {square.id}
        </div>
    )
}

const squareStyle = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
}

ChessSquare.defaultProps = {
    square: {},
    clickSquare: () => {}
}

ChessSquare.propTypes = {
    square: PropTypes.object,
    clickSquare: PropTypes.func
}

export default ChessSquare
