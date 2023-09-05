import React from 'react'

const pieceStyle: React.CSSProperties = {
    width: '100%',
    height: '100%',
}

const ChessPiece = ({piece}: {piece: any}): JSX.Element => {
    return (
        <div
            style={{...pieceStyle, color: piece.color}}
        >
            {piece.constructor.name}
        </div>
    )
}

export default ChessPiece