import { useState } from 'react'
import ChessBoard from './components/ChessBoard.jsx'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  const boardConfig = {
    xSquares: 8,
    ySquares: 8,
    startingPlayer: 'white',
    playerColor: 'white',
    startingPieces: {
      '00': {name: 'r', color: 'black', moved: false},
      '10': {name: 'n', color: 'black'},
      '20': {name: 'b', color: 'black'},
      '30': {name: 'q', color: 'black'},
      '40': {name: 'k', color: 'black', moved: false, inCheck: false, inCheckmate: false},
      '50': {name: 'b', color: 'black'},
      '60': {name: 'n', color: 'black'},
      '70': {name: 'r', color: 'black', moved: false},
      '01': {name: 'p', color: 'black', moved: false, enPassant: '', xDir: 0, yDir: 1},
      '11': {name: 'p', color: 'black', moved: false, enPassant: '', xDir: 0, yDir: 1},
      '21': {name: 'p', color: 'black', moved: false, enPassant: '', xDir: 0, yDir: 1},
      '31': {name: 'p', color: 'black', moved: false, enPassant: '', xDir: 0, yDir: 1},
      '41': {name: 'p', color: 'black', moved: false, enPassant: '', xDir: 0, yDir: 1},
      '51': {name: 'p', color: 'black', moved: false, enPassant: '', xDir: 0, yDir: 1},
      '61': {name: 'p', color: 'black', moved: false, enPassant: '', xDir: 0, yDir: 1},
      '71': {name: 'p', color: 'black', moved: false, enPassant: '', xDir: 0, yDir: 1},
      '06': {name: 'p', color: 'white', moved: false, enPassant: '', xDir: 0, yDir: -1},
      '16': {name: 'p', color: 'white', moved: false, enPassant: '', xDir: 0, yDir: -1},
      '26': {name: 'p', color: 'white', moved: false, enPassant: '', xDir: 0, yDir: -1},
      '36': {name: 'p', color: 'white', moved: false, enPassant: '', xDir: 0, yDir: -1},
      '46': {name: 'p', color: 'white', moved: false, enPassant: '', xDir: 0, yDir: -1},
      '56': {name: 'p', color: 'white', moved: false, enPassant: '', xDir: 0, yDir: -1},
      '66': {name: 'p', color: 'white', moved: false, enPassant: '', xDir: 0, yDir: -1},
      '76': {name: 'p', color: 'white', moved: false, enPassant: '', xDir: 0, yDir: -1},
      '07': {name: 'r', color: 'white', moved: false},
      '17': {name: 'n', color: 'white'},
      '27': {name: 'b', color: 'white'},
      '37': {name: 'q', color: 'white'},
      '47': {name: 'k', color: 'white', moved: false, inCheck: false, inCheckmate: false},
      '57': {name: 'b', color: 'white'},
      '67': {name: 'n', color: 'white'},
      '77': {name: 'r', color: 'white', moved: false},
    },
  }

  return (
    <>
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.jsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
      <ChessBoard boardConfig={boardConfig}/>
    </>
  )
}

export default App
