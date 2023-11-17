import { useState, useEffect } from 'react'
import useWebSocket from 'react-use-websocket'
import ChessBoard from './components/ChessBoard'


const CreateMoveData = (xFrom: number, yFrom: number, xTo: number, yTo: number) : string => {
    return JSON.stringify({
        title: 'move',
        data: {
            xFrom: xFrom,
            yFrom: yFrom,
            xTo: xTo,
            yTo: yTo,
        }
    })
}

const CreateViewData = (x: number, y: number) : string => {
    return JSON.stringify({
        title: 'view',
        data: {
            x: x,
            y: y,
        }
    })
}

const createJsonMessage = (message: string) : string => {
    return JSON.stringify({
        type: 'message',
        data: message
    })
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
    const [message, setMessage] = useState('')
    const [socketUrl] = useState(`ws://localhost:8080${urlExtension}`)
    const { sendMessage, lastMessage, readyState, getWebSocket } = useWebSocket(socketUrl)
    
    useEffect(() => {
        console.log(lastMessage)
    }, [lastMessage])

    // setTimeout(() => {
    //     getWebSocket()?.close()
    // }, 5000)

    const handleMove = (xFrom: number, yFrom: number, xTo: number, yTo: number): void => {
        const moveData = CreateMoveData(xFrom, yFrom, xTo, yTo)
        sendMessage(moveData)
    }

    const handleView = (x: number, y: number): void => {
        const viewData = CreateViewData(x, y)
        sendMessage(viewData)
    }

    return (
        <>
            <input
                type="text"
                onChange={(e) => setMessage(e.target.value)}
            />
            <button
                onClick={() => sendMessage(createJsonMessage(message))}
            />
            <ChessBoard handleMove={handleMove} handleView={handleView} />
        </>
    )
}

export default ClientTest
