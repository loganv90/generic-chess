import React from 'react'

const actionBarStyle: React.CSSProperties = {
    width: '500px',
    height: '100px',
}

const ChessActionBar = ({undo, redo}: {undo: () => void, redo: () => void}): JSX.Element => {
    return (
        <div style={actionBarStyle}>
            <button onClick={undo}>Undo</button>
            <button onClick={redo}>Redo</button>
        </div>
    )
}

export default ChessActionBar