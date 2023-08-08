import PropTypes from 'prop-types'

function ChessPiece({ piece }) {
    return (
        <div
            style={{...pieceStyle, color: piece.color}}
        >
            {piece.constructor.name}
        </div>
    )
}

const pieceStyle = {
    width: '100%',
    height: '100%',
}

ChessPiece.defaultProps = {
    piece: {},
    clickSquare: () => {}
}

ChessPiece.propTypes = {
    piece: PropTypes.object,
    clickPiece: PropTypes.func
}

export default ChessPiece
