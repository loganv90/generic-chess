import React from 'react'

const pieceStyle: React.CSSProperties = {
    width: '100%',
    height: '100%',
}

const ChessPiece = ({
    type,
    color,
}: {
    type: string,
    color: string,
}): JSX.Element => {
    const getPieceColor = (): string => {
        if (color === 'white') {
            return 'white'
        } else {
            return 'black'
        }
    }

    return (
        <div
            style={{...pieceStyle, color: getPieceColor()}}
        >
            {type}
        </div>
    )
}

export default ChessPiece
