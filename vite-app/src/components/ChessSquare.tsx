import React from 'react'
import ChessPiece from './ChessPiece'
import { SquareData } from './ChessBoard'

const squareStyle: React.CSSProperties = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
}

const ChessSquare = ({
    sizePercent,
    square,
    selected,
    destination,
    clickSquare
}:{
    sizePercent: string,
    square: SquareData,
    selected: boolean,
    destination: boolean,
    clickSquare: (id: string, x: number, y: number) => void
}): JSX.Element => {
    const getSquareColor = (): string => {
        if (selected) {
            return 'red'
        } else if (destination) {
            return 'blue'
        } else if (square.light) {
            return '#306010'
        } else {
            return '#304030'
        }
    }

    return (
        <div
            style={{...squareStyle, backgroundColor: getSquareColor(), width: sizePercent, height: sizePercent}}
            onClick={() => clickSquare(square.id, square.x, square.y)}
        >
            <ChessPiece type={square.type} color={square.color} />
        </div>
    )
}

export default ChessSquare
