import { useState, useEffect, useRef } from 'react'
import './App.css'

function App() {
  const [gameState, setGameState] = useState(null)
  const [status, setStatus] = useState("Connecting to server...")
  const ws = useRef(null)

  useEffect(() => {
    // Connect to your Go Backend
    const WS_URL = "wss://four-in-a-row-ilyp.onrender.com/ws"; 

ws.current = new WebSocket(WS_URL);

    ws.current = new WebSocket("ws://localhost:8080/ws")

    ws.current.onopen = () => {
      setStatus("Searching for opponent... (Wait 10s for Bot)")
    }

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data)
      setGameState(data)
      
      if (data.winner > 0) {
        if (data.winner === 3) setStatus("It's a Draw!")
        else setStatus(data.winner === data.youAre ? "🎉 You Won!" : "💀 You Lost!")
      } else {
        setStatus(data.turn === data.youAre ? "🟢 Your Turn" : "🔴 Opponent's Turn")
      }
    }

    return () => ws.current.close()
  }, [])

  const handleDrop = (colIndex) => {
    if (!gameState || gameState.winner > 0) return
    if (gameState.turn !== gameState.youAre) return // Not your turn

    ws.current.send(JSON.stringify({ col: colIndex }))
  }

  // Render Board with drop indicator
  const renderBoard = () => {
    if (!gameState) return null
    return (
      <>
        {/* Drop indicator row */}
        <div className="drop-row">
          {Array(7).fill(0).map((_, cIndex) => (
            <div
              key={`drop-${cIndex}`}
              className="drop-indicator"
              onClick={() => handleDrop(cIndex)}
            >
              <span role="img" aria-label="drop">⬇️</span>
            </div>
          ))}
        </div>
        {/* Game board grid */}
        <div className="grid">
          {gameState.board.map((row, rIndex) =>
            row.map((cell, cIndex) => {
              let classColor = ""
              if (cell === 1) classColor = "red"
              if (cell === 2) classColor = "yellow"
              return (
                <div
                  key={`${rIndex}-${cIndex}`}
                  className={`cell ${classColor}`}
                  style={{ animationDelay: `${rIndex * 0.05 + cIndex * 0.02}s` }}
                  onClick={() => handleDrop(cIndex)}
                />
              )
            })
          )}
        </div>
      </>
    )
  }

  return (
    <div className="game-container">
      <h1 className="game-title">Connect 4 <span className="realtime-badge">Real Time</span></h1>
      <div className="status">{status}</div>
      <div className="board">
        {gameState ? renderBoard() : <p>Loading Game...</p>}
      </div>
    </div>
  )
}

export default App