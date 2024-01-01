import React from 'react'
import ChessPiece from './ChessPiece'
import { SquareData } from './ChessBoard'

const squareStyle: React.CSSProperties = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    width: '12.5%',
    height: '12.5%',
}

const ChessSquare = ({
    square,
    selected,
    destination,
    clickSquare
}:{
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
            return '#304030'
        } else {
            return '#306010'
        }
    }

    return (
        <div
            style={{...squareStyle, backgroundColor: getSquareColor()}}
            onClick={() => clickSquare(square.id, square.x, square.y)}
        >
            <ChessPiece type={square.type} color={square.color} />
        </div>
    )
}

export default ChessSquare
