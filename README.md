# 4 in a Row - Realtime Multiplayer Game

A real-time, backend-driven Connect 4 game built with **Go (Golang)** and **React**. Supports 1v1 PvP and intelligent Bot fallback.

**[Link to Live App](PUT_YOUR_VERCEL_LINK_HERE)**

## 🚀 Features
* **Real-time Multiplayer:** Uses WebSockets for instant state synchronization.
* **Smart Matchmaking:** * Wait 10s for a human opponent.
    * If no match found, seamlessly switch to an AI Bot.
* **Resiliency:** Handles ghost connections and disconnections gracefully.
* **Game Logic:** Full server-side validation (the server is the "source of truth").

## 🛠️ Tech Stack
* **Backend:** Go 1.25, Gorilla WebSocket, Goroutines (Concurrency)
* **Frontend:** React, Vite, CSS Grid
* **Deployment:** Render (Backend), Vercel (Frontend)

## 🏃‍♂️ How to Run Locally

### Prerequisites
* Go 1.18+
* Node.js & npm

### 1. Start the Backend
```bash
cd backend
go mod tidy
go run cmd/server/main.go
# Server starts on localhost:8080
