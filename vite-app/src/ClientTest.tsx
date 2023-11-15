import { useState, useEffect } from 'react'
import useWebSocket from 'react-use-websocket'
import ChessBoard from './components/ChessBoard'

function createJsonMessage(message: string) : string {
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

    return (
        <>
            <input
                type="text"
                onChange={(e) => setMessage(e.target.value)}
            />
            <button
                onClick={() => sendMessage(createJsonMessage(message))}
            />
            <button
                onClick={() => sendMessage(CreateMoveData())}
            >move data button</button>
            <ChessBoard />
        </>
    )
}

const CreateMoveData = () : string => {
    return JSON.stringify({
        title: 'move',
        data: {
            xFrom: 4,
            yFrom: 6,
            xTo: 4,
            yTo: 4,
        }
    })
}

export default ClientTest
