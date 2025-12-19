package game

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/abhisuman2001/connect4/internal/models"
)

type Room struct {
	ID        string
	Logic     *GameLogic
	Player1   *models.Client
	Player2   *models.Client // Nil if bot
	IsBot     bool
	Mutex     sync.Mutex
}

type Manager struct {
	Queue chan *models.Client
	Rooms map[string]*Room
	Mutex sync.RWMutex
}

func NewManager() *Manager {
	m := &Manager{
		Queue: make(chan *models.Client, 100),
		Rooms: make(map[string]*Room),
	}
	go m.matchmaker()
	return m
}

func (m *Manager) matchmaker() {
	for {
		// 1. Pick the first player
		p1 := <-m.Queue
		log.Println("Player 1 joined queue, waiting...")

		timer := time.NewTimer(10 * time.Second)

		select {
		case p2 := <-m.Queue:
			timer.Stop()
			
			// 🔍 CHECK: Is Player 1 still there?
			// We try to send a "Ping". If it fails, P1 is a ghost.
			if err := p1.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Player 1 disconnected/ghosted. Starting Bot Game for Player 2.")
				// P1 is dead. We match P2 with a Bot immediately so they don't get stuck.
				m.CreateRoom(p2, nil, true)
			} else {
				// Both are alive! Start PvP.
				m.CreateRoom(p1, p2, false)
			}

		case <-timer.C:
			log.Println("Timeout: Starting Bot Game")
			m.CreateRoom(p1, nil, true)
		}
	}
}

func (m *Manager) CreateRoom(p1, p2 *models.Client, isBot bool) {
	roomID := uuid.New().String()
	p1.GameID = roomID
	p1.Color = models.Red

	room := &Room{
		ID:      roomID,
		Logic:   NewGameLogic(),
		Player1: p1,
		IsBot:   isBot,
	}

	if !isBot {
		p2.GameID = roomID
		p2.Color = models.Yellow
		room.Player2 = p2
	}

	m.Mutex.Lock()
	m.Rooms[roomID] = room
	m.Mutex.Unlock()

	room.BroadcastState()
}

func (r *Room) HandleMove(client *models.Client, col int) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.Logic.DropPiece(col, client.Color) {
		r.BroadcastState()
		
		if r.IsBot && !r.Logic.GameOver {
			time.Sleep(500 * time.Millisecond) // Fake think time
			botCol := r.Logic.GetBotMove()
			if botCol != -1 {
				r.Logic.DropPiece(botCol, models.Yellow)
				r.BroadcastState()
			}
		}
	}
}

func (r *Room) BroadcastState() {
	baseState := models.GameState{
		Type:   "update",
		Board:  r.Logic.Board,
		Turn:   r.Logic.Turn,
		Winner: r.Logic.Winner,
	}

	// Send to Player 1
	baseState.YouAre = models.Red
	sendJSON(r.Player1.Conn, baseState)

	// Send to Player 2
	if !r.IsBot && r.Player2 != nil {
		baseState.YouAre = models.Yellow
		sendJSON(r.Player2.Conn, baseState)
	}
}

func sendJSON(conn *websocket.Conn, v interface{}) {
	msg, _ := json.Marshal(v)
	conn.WriteMessage(websocket.TextMessage, msg)
}