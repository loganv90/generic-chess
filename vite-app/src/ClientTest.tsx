import { useState, useEffect } from 'react'
import useWebSocket from 'react-use-websocket'

function ClientTest() {
    const [socketUrl] = useState('ws://localhost:8080/ws')
    const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl)
    
    useEffect(() => {
        console.log(lastMessage)
    }, [lastMessage])

    return (
        <>
            <div>
                {readyState}
            </div>
        </>
    )
}

export default ClientTest
