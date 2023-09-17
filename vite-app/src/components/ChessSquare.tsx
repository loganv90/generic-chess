import React from 'react'
import ChessPiece from './ChessPiece'
import { Square } from '../utils/chess/Square'
import { Piece } from '../utils/chess/Piece'

const squareStyle: React.CSSProperties = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    width: '12.5%',
    height: '12.5%',
}

const ChessSquare = ({
    square,
    isSelected,
    isDestination,
    clickSquare
}:{
    square: Square,
    isSelected: boolean,
    isDestination: boolean,
    clickSquare: (x: number, y: number) => void
}): JSX.Element => {
    const piece: Piece | null = square.getPiece()

    const getSquareColor = (): string => {
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
            style={{...squareStyle, backgroundColor: getSquareColor()}}
            onClick={() => clickSquare(square.x, square.y)}
        >
            {piece && <ChessPiece piece={piece} />}
        </div>
    )
}

export default ChessSquare