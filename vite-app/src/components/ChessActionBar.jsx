import PropTypes from 'prop-types'

function ChessActionBar({undo, redo}) {
    return (
        <div style={{...actionBarStyle}}>
            <button onClick={undo}>Undo</button>
            <button onClick={redo}>Redo</button>
        </div>
    )
}

const actionBarStyle = {
    width: '500px',
    height: '100px',
}

ChessActionBar.defaultProps = {
    undo: () => {},
    redo: () => {},
}

ChessActionBar.propTypes = {
    undo: PropTypes.func,
    redo: PropTypes.func,
}

export default ChessActionBar
