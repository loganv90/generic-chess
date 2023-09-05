import React from 'react'
import ChessPiece from './ChessPiece'

const squareStyle: React.CSSProperties = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    width: '12.5%',
    height: '12.5%',
}

const ChessSquare = (
    { square, isSelected, isDestination, clickSquare }:
    { square: any, isSelected: boolean, isDestination: boolean, clickSquare: Function }
): JSX.Element => {
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
            {square.piece && <ChessPiece piece={square.piece} />}
        </div>
    )
}

export default ChessSquare