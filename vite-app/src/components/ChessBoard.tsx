import React, { useState, useEffect } from 'react'
import ChessSquare from './ChessSquare'
import ChessActionBar from './ChessActionBar'
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

const createDestinations = (moveData: MoveData): Set<string> => {
    const destinations = new Set<string>()

    moveData.Moves.forEach(m => destinations.add(`${m.X},${m.Y}`))

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
    undo,
    redo,
}: {
    boardData: BoardData,
    moveData: MoveData,
    move: (xFrom: number, yFrom: number, xTo: number, yTo: number) => void,
    view: (x: number, y: number) => void,
    undo: () => void,
    redo: () => void,
}): JSX.Element => {
    const [squares, setSquares] = useState<SquareData[][]>(createSquares(boardData))
    const [destinations, setDestinations] = useState<Set<string>>(createDestinations(moveData))
    const [selected, setSelected] = useState<string>(createSelected(moveData))

    useEffect(() => {
        setSquares(createSquares(boardData))
    }, [boardData])

    useEffect(() => {
        setDestinations(createDestinations(moveData))
        setSelected(createSelected(moveData))
    }, [moveData])

    const undoMove = (): void => {
        setSquares(s => [...s])
        undo()
    }

    const redoMove = (): void => {
        setSquares(s => [...s])
        redo()
    }

    const clickSquare = (id: string, x: number, y: number): void => {
        if (selected === id) {
            setSelected('-1,-1')
            setDestinations(new Set<string>())
            return
        }

        if (destinations.has(id)) {
            move(moveData.X, moveData.Y, x, y)
            setSelected('-1,-1')
            setDestinations(new Set<string>())
            return
        }

        view(x, y)
    }

    return (
        <div>
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
            <ChessActionBar undo={undoMove} redo={redoMove} />
        </div>
    )
}

export default ChessBoard
export type { SquareData }
