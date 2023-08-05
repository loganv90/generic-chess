import PropTypes from 'prop-types'
import ChessPiece from './ChessPiece'

function ChessSquare({ square, isSelected, isDestination, clickSquare,  }) {
    const getSquareColor = () => {
        if (isSelected) {
            return 'red'
        } else if (isDestination) {
            return 'blue'
        } else if (square.isLight) {
            return 'lightGrey'
        } else {
            return 'darkGrey'
        }
    }

    return (
        <div
            style={{
                ...squareStyle,
                backgroundColor: getSquareColor(),
            }}
            onClick={() => clickSquare(square.x, square.y)}
        >
            {square.piece && <ChessPiece pieceConfig={square.piece} />}
        </div>
    )
}

const squareStyle = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    width: '12.5%',
    height: '12.5%',
}

ChessSquare.defaultProps = {
    square: {},
    isSelected: false,
    isDestination: false,
    clickSquare: () => {}
}

ChessSquare.propTypes = {
    square: PropTypes.object,
    isSelected: PropTypes.bool,
    isDestination: PropTypes.bool,
    clickSquare: PropTypes.func
}

export default ChessSquare
