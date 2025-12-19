package models

import "github.com/gorilla/websocket"

const (
	Empty  = 0
	Red    = 1 // Player 1
	Yellow = 2 // Player 2
	Rows   = 6
	Cols   = 7
)

// Message sent to frontend
type GameState struct {
	Type     string         `json:"type"` // "start", "update", "end"
	Board    [Rows][Cols]int `json:"board"`
	Turn     int            `json:"turn"`
	Winner   int            `json:"winner"` // 0=none, 1=Red, 2=Yellow, 3=Draw
	YouAre   int            `json:"youAre"` // 1 or 2 (assigned to client)
	Message  string         `json:"message"`
}

// Message received from frontend
type PlayerMove struct {
	Col int `json:"col"`
}

type Client struct {
	Conn   *websocket.Conn
	Color  int
	GameID string
}