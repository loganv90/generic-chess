import PropTypes from 'prop-types'
import ChessPiece from './ChessPiece'

function ChessSquare({ squareConfig, clickSquare }) {
    const getSquareColor = () => {
        if (squareConfig.isSelected) {
            return 'red'
        } else if (squareConfig.isDark) {
            return 'darkGrey'
        } else {
            return 'lightGrey'
        }
    }

    const clickPiece = () => {
        clickSquare(squareConfig.id)
    }

    return (
        <div
            style={{
                ...squareStyle,
                backgroundColor: getSquareColor(),
                width: squareConfig.width,
                height: squareConfig.height
            }}
        >
            {squareConfig.piece && <ChessPiece pieceConfig={squareConfig.piece} clickPiece={clickPiece} />}
        </div>
    )
}

const squareStyle = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
}

ChessSquare.defaultProps = {
    squareConfig: {},
    clickSquare: () => {}
}

ChessSquare.propTypes = {
    squareConfig: PropTypes.object,
    clickSquare: PropTypes.func
}

export default ChessSquare
