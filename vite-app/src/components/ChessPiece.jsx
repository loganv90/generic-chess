import PropTypes from 'prop-types'

function ChessPiece({ pieceConfig, clickPiece }) {
    return (
        <div
            style={{...pieceStyle, color: pieceConfig.color}}
            onClick={() => clickPiece()}
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
