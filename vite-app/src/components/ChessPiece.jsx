import PropTypes from 'prop-types'

function ChessPiece({ pieceConfig }) {
    return (
        <div
            style={{...pieceStyle, color: pieceConfig.color}}
        >
            {pieceConfig.name}
        </div>
    )
}

const pieceStyle = {
    width: '100%',
    height: '100%',
}

ChessPiece.defaultProps = {
    pieceConfig: {},
    clickSquare: () => {}
}

ChessPiece.propTypes = {
    pieceConfig: PropTypes.object,
    clickPiece: PropTypes.func
}

export default ChessPiece
