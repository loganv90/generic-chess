import React from 'react'
import { Piece } from '../utils/chess/Piece'

const pieceStyle: React.CSSProperties = {
    width: '100%',
    height: '100%',
}

const ChessPiece = ({piece}: {piece: Piece}): JSX.Element => {
    return (
        <div
            style={{...pieceStyle, color: piece.color}}
        >
            {piece.constructor.name}
        </div>
    )
}

export default ChessPiece