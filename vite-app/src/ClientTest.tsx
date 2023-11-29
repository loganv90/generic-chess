import { useState, useEffect } from 'react'
import useWebSocket from 'react-use-websocket'
import ChessBoard from './components/ChessBoard'

type Message = {
    Type: string,
    Data: any
}

const createMoveData = (xFrom: number, yFrom: number, xTo: number, yTo: number) : Message => {
    return {
        Type: 'move',
        Data: {
            xFrom: xFrom,
            yFrom: yFrom,
            xTo: xTo,
            yTo: yTo,
        }
    }
}

const createViewData = (x: number, y: number) : Message => {
    return {
        Type: 'view',
        Data: {
            x: x,
            y: y,
        }
    }
}

const baseUrl = '/ws'
const ClientTest = () => {
    const [urlExtension, setUrlExtension] = useState('')
    const [currentUrlExtension, setCurrentUrlExtension] = useState('')

    return (
        <>
            <input
                type="text"
                onChange={(e) => setUrlExtension(e.target.value)}
            />
            <button
                onClick={() => setCurrentUrlExtension(`${baseUrl}${urlExtension}`)}
            />
            { currentUrlExtension && <ChessGame urlExtension={currentUrlExtension} /> }
        </>
    )
}


const ChessGame = ({ urlExtension }: { urlExtension: string }): JSX.Element => {
    const [socketUrl] = useState(`ws://localhost:8080${urlExtension}`)
    const { sendMessage, lastMessage, readyState, getWebSocket } = useWebSocket(socketUrl)
    // setTimeout(() => {
    //     getWebSocket()?.close()
    // }, 5000)
    
    useEffect(() => {
        console.log(lastMessage)
    }, [lastMessage])

    const handleMove = (xFrom: number, yFrom: number, xTo: number, yTo: number): void => {
        const moveData = createMoveData(xFrom, yFrom, xTo, yTo)
        sendMessage(JSON.stringify(moveData))
    }

    const handleView = (x: number, y: number): void => {
        const viewData = createViewData(x, y)
        sendMessage(JSON.stringify(viewData))
    }

    const handleUndo = (): void => {
        console.log('undo')
    }

    const handleRedo = (): void => {
        console.log('redo')
    }

    return (
        <ChessBoard
            handleMove={handleMove}
            handleView={handleView}
            handleUndo={handleUndo}
            handleRedo={handleRedo}
        />
    )
}

export default ClientTest
