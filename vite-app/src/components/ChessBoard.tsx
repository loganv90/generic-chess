import React, { useState, useEffect } from 'react'
import ChessSquare from './ChessSquare'
import { BoardData, MoveData } from '../ClientTest.tsx'

const boardStyle: React.CSSProperties = {
    display: 'flex',
    flexWrap: 'wrap',
    width: '500px',
    height: '500px',
}

type SquareData = {
    x: number,
    y: number,
    id: string,
    light: boolean,
    type: string,
    color: string,
}

const createSquares = (boardData: BoardData): SquareData[][] => {
    const squares: SquareData[][] = []
    
    for (let y = 0; y < boardData.YSize; y++) {
        squares.push([])
        for (let x = 0; x < boardData.XSize; x++) {
            squares[y].push({
                x: x,
                y: y,
                id: `${x},${y}`,
                light: (x + y) % 2 === 0,
                type: '',
                color: '',
            })
        }
    }

    boardData.Pieces.forEach(p => {
        squares[p.Y][p.X].type = p.T
        squares[p.Y][p.X].color = p.C
    })

    return squares
}

const createDestinations = (moveData: MoveData): Map<string, boolean> => {
    const destinations = new Map<string, boolean>()

    moveData.Moves.forEach(m => destinations.set(`${m.X},${m.Y}`, m.P))

    return destinations
}

const createSelected = (moveData: MoveData): string => {
    return `${moveData.X},${moveData.Y}`
}

const ChessBoard = ({
    boardData,
    moveData,
    move,
    view,
}: {
    boardData: BoardData,
    moveData: MoveData,
    move: (xFrom: number, yFrom: number, xTo: number, yTo: number) => void,
    view: (x: number, y: number) => void,
}): JSX.Element => {
    const [squares, setSquares] = useState<SquareData[][]>(createSquares(boardData))
    const [destinations, setDestinations] = useState<Map<string, boolean>>(createDestinations(moveData))
    const [selected, setSelected] = useState<string>(createSelected(moveData))

    useEffect(() => {
        setSquares(createSquares(boardData))
    }, [boardData])

    useEffect(() => {
        setDestinations(createDestinations(moveData))
        setSelected(createSelected(moveData))
    }, [moveData])

    const clickSquare = (id: string, x: number, y: number): void => {
        if (selected === id) {
            setSelected('-1,-1')
            setDestinations(new Map<string, boolean>())
            return
        }

        if (!destinations.has(id)) {
            view(x, y)
            return
        }

        if (destinations.get(id)) {
            move(moveData.X, moveData.Y, x, y)
        } else {
            move(moveData.X, moveData.Y, x, y)
        }

        setSelected('-1,-1')
        setDestinations(new Map<string, boolean>())
        return
    }

    return (
        <div style={boardStyle}>
            {squares.map((row) => row.map((s) => 
                <ChessSquare
                    key={s.id}
                    square={s}
                    selected={selected === s.id}
                    destination={destinations.has(s.id)}
                    clickSquare={clickSquare}
                />
            ))}
        </div>
    )
}

export default ChessBoard
export type { SquareData }
