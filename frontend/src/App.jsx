import { useState, useEffect, useRef } from 'react'
import './App.css'

function App() {
  const [gameState, setGameState] = useState(null)
  const [status, setStatus] = useState("Connecting...")
  const [username, setUsername] = useState("")
  const [isLoggedIn, setIsLoggedIn] = useState(false) // <--- NEW STATE
  const ws = useRef(null)

  useEffect(() => {
    // Only connect IF user has logged in
    if (!isLoggedIn) return;

    // Connect to Backend
    const WS_URL = "wss://four-in-a-row-ilyp.onrender.com/ws";
    ws.current = new WebSocket(WS_URL);

    ws.current.onopen = () => {
      setStatus("Searching for opponent... (Wait 10s for Bot)")
    }

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data)
      setGameState(data)
      
      if (data.winner > 0) {
        if (data.winner === 3) setStatus("It's a Draw!")
        else setStatus(data.winner === data.youAre ? `🎉 ${username}, You Won!` : "💀 You Lost!")
      } else {
        setStatus(data.turn === data.youAre ? `🟢 Your Turn, ${username}` : "🔴 Opponent's Turn")
      }
    }

    return () => {
      if (ws.current) ws.current.close()
    }
  }, [isLoggedIn]) // <--- Depends on isLoggedIn

  const handleDrop = (colIndex) => {
    if (!gameState || gameState.winner > 0) return
    if (gameState.turn !== gameState.youAre) return

    ws.current.send(JSON.stringify({ col: colIndex }))
  }

  const handleLogin = () => {
    if (username.trim() !== "") {
      setIsLoggedIn(true);
    }
  }

  const renderBoard = () => {
    if (!gameState) return null
    return gameState.board.map((row, rIndex) => 
      row.map((cell, cIndex) => {
        let classColor = ""
        if (cell === 1) classColor = "red"
        if (cell === 2) classColor = "yellow"
        
        return (
          <div 
            key={`${rIndex}-${cIndex}`} 
            className={`cell ${classColor}`}
            onClick={() => handleDrop(cIndex)}
          />
        )
      })
    )
  }

  // --- 1. RENDER LOGIN SCREEN IF NOT LOGGED IN ---
  if (!isLoggedIn) {
    return (
      <div className="login-screen">
        <h1>4 in a Row</h1>
        <p>Enter your username to join the battle!</p>
        <input 
          type="text" 
          placeholder="Username..." 
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          style={{ padding: '10px', fontSize: '16px', borderRadius: '5px', border: '1px solid #ccc' }}
        />
        <br />
        <button 
          onClick={handleLogin}
          style={{ marginTop: '20px', padding: '10px 20px', fontSize: '16px', background: '#3b82f6', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer' }}
        >
          Join Game
        </button>
      </div>
    )
  }

  // --- 2. RENDER GAME BOARD ---
  return (
    <div>
      <h1>Connect 4 - Real Time</h1>
      <div className="status">{status}</div>
      {gameState && (
        <div className="board">
          {renderBoard()}
        </div>
      )}
      {!gameState && <p>Loading Game Board...</p>}
    </div>
  )
}

export default App