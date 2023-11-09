import { useState, useEffect } from 'react'
import useWebSocket from 'react-use-websocket'

function ClientTest() {
    const [socketUrl] = useState('ws://localhost:8080/ws')
    const { sendMessage, lastMessage, readyState, getWebSocket } = useWebSocket(socketUrl)
    
    useEffect(() => {
        console.log(lastMessage)
    }, [lastMessage])

    setTimeout(() => {
        // close websocket connection
        getWebSocket()?.close()
    }, 5000)

    return (
        <>
            <div>
                {readyState}
            </div>
        </>
    )
}

export default ClientTest
