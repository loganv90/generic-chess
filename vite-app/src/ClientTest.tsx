import { useState, useEffect } from 'react'
import useWebSocket from 'react-use-websocket'
import ChessBoard from './components/ChessBoard'

type Message = {
    Type: string,
    Data: any,
}

type BoardData = {
    XSize: number,
    YSize: number,
    Disabled: { X: number, Y: number }[],
    Pieces: { T: string, C: string, X: number, Y: number }[],
    CurrentPlayer: string,
    WinningPlayer: string,
    Check: boolean,
    Checkmate: boolean,
    Stalemate: boolean,
}

type MoveData = {
    X: number,
    Y: number,
    Moves: { X: number, Y: number, P: boolean }[],
}

const createMoveMessage = (xFrom: number, yFrom: number, xTo: number, yTo: number, promotion: string) : Message => {
    return {
        Type: 'move',
        Data: {
            xFrom: xFrom,
            yFrom: yFrom,
            xTo: xTo,
            yTo: yTo,
            promotion: promotion,
        }
    }
}

const createViewMessage = (x: number, y: number) : Message => {
    return {
        Type: 'view',
        Data: {
            x: x,
            y: y,
        }
    }
}

const createUndoMessage = () : Message => {
    return {
        Type: 'undo',
        Data: {}
    }
}

const createRedoMessage = () : Message => {
    return {
        Type: 'redo',
        Data: {}
    }
}

const baseUrl = '/ws'
const ClientTest = () => {
    const [urlExtension, setUrlExtension] = useState('')
    const [currentUrlExtension, setCurrentUrlExtension] = useState('')

    return (
        <>
            <div style={{textAlign: "left"}}>
                <label htmlFor="two">Two player: </label>
                <button id="two" name="two" onClick={() => setCurrentUrlExtension(`${baseUrl}/two`)} />
                <br />
                <label htmlFor="four">Four player: </label>
                <button id="four" name="four" onClick={() => setCurrentUrlExtension(`${baseUrl}/four`)} />
                <br />
                <label htmlFor="code">Game code: </label>
                <input type="text" id="code" name="code" onChange={(e) => setUrlExtension(e.target.value)} />
                <button id="code" name="code" onClick={() => setCurrentUrlExtension(`${baseUrl}/${urlExtension}`)} />
            </div>
            <br />
            { currentUrlExtension && <ChessGame urlExtension={currentUrlExtension} /> }
        </>
    )
}

const ChessGame = ({ urlExtension }: { urlExtension: string }): JSX.Element => {
    const [socketUrl] = useState(`ws://localhost:8080${urlExtension}`)
    const { sendMessage, lastMessage, readyState, getWebSocket } = useWebSocket(socketUrl)
    const [promotionPiece, setPromotionPiece] = useState('Q')

    const [boardData, setBoardData] = useState<BoardData>({
        XSize: 0,
        YSize: 0,
        Disabled: [],
        Pieces: [],
        CurrentPlayer: '',
        WinningPlayer: '',
        Check: false,
        Checkmate: false,
        Stalemate: false,
    })

    const [moveData, setMoveData] = useState<MoveData>({
        X: -1,
        Y: -1,
        Moves: [],
    })

    useEffect(() => {
        const message = JSON.parse(lastMessage?.data || '{}') as Message

        if (message.Type === 'BoardState') {
            const boardData = message.Data as BoardData
            console.log(boardData)
            setBoardData(boardData)
        } else if (message.Type === 'PieceState') {
            const moveData = message.Data as MoveData
            console.log(moveData)
            setMoveData(moveData)
        } else {
            console.log('unknown message', lastMessage)
        }
    }, [lastMessage])

    const handleMove = (xFrom: number, yFrom: number, xTo: number, yTo: number): void => {
        const moveMessage = createMoveMessage(xFrom, yFrom, xTo, yTo, promotionPiece)
        sendMessage(JSON.stringify(moveMessage))
    }

    const handleView = (x: number, y: number): void => {
        const viewMessage = createViewMessage(x, y)
        sendMessage(JSON.stringify(viewMessage))
    }

    const handleUndo = (): void => {
        const undoMessage = createUndoMessage()
        sendMessage(JSON.stringify(undoMessage))
    }

    const handleRedo = (): void => {
        const redoMessage = createRedoMessage()
        sendMessage(JSON.stringify(redoMessage))
    }

    const handlePromotionPiece = (e: React.ChangeEvent<HTMLSelectElement>): void => {
        setPromotionPiece(e.target.value)
    }

    return (
        <>
            <ChessBoard
                boardData={boardData}
                moveData={moveData}
                move={handleMove}
                view={handleView}
            />
            <div style={{textAlign: "left"}}>
                <br />
                Current player: {boardData.CurrentPlayer}
                <br />
                Winning player: {boardData.WinningPlayer}
                <br />
                <label htmlFor="pieces">Default promotion piece: </label>
                <select id="pieces" name="pieces" value={promotionPiece} onChange={handlePromotionPiece} >
                  <option value="Q">Queen</option>
                  <option value="R">Rook</option>
                  <option value="B">Bishop</option>
                  <option value="N">Knight</option>
                </select>
                <br />
                <label htmlFor="undo">Undo move: </label>
                <button id="undo" name="undo" onClick={handleUndo} />
                <br />
                <label htmlFor="redo">Redo move: </label>
                <button id="redo" name="redo" onClick={handleRedo} />
                <br />
            </div>
        </>
    )
}

export default ClientTest
export type { BoardData, MoveData }
