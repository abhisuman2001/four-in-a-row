package main

import (
	"os"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/abhisuman2001/connect4/internal/game"
	"github.com/abhisuman2001/connect4/internal/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var gameManager = game.NewManager()

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	client := &models.Client{Conn: conn}
	gameManager.Queue <- client

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		
		var move models.PlayerMove
		if err := json.Unmarshal(msg, &move); err == nil {
			gameManager.Mutex.RLock()
			room, ok := gameManager.Rooms[client.GameID]
			gameManager.Mutex.RUnlock()
			
			if ok {
				room.HandleMove(client, move.Col)
			}
		}
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("✅ 4-in-a-Row Game Server is Running!"))
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWS)
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}