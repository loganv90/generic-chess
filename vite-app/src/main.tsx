import React from 'react'
import ReactDOM from 'react-dom/client'
import ChessClient from './components/ChessClient'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <ChessClient />
    </React.StrictMode>,
)

