package game

import (
	"math/rand"
	"github.com/abhisuman2001/connect4/internal/models"
)

type GameLogic struct {
	Board    [models.Rows][models.Cols]int
	Turn     int
	Winner   int
	GameOver bool
	Moves    int
}

func NewGameLogic() *GameLogic {
	return &GameLogic{
		Turn: models.Red,
	}
}

func (g *GameLogic) DropPiece(col int, player int) bool {
	if col < 0 || col >= models.Cols || g.GameOver || player != g.Turn {
		return false
	}

	for r := models.Rows - 1; r >= 0; r-- {
		if g.Board[r][col] == models.Empty {
			g.Board[r][col] = player
			g.Moves++
			
			if g.checkWin(r, col, player) {
				g.Winner = player
				g.GameOver = true
			} else if g.Moves >= models.Rows*models.Cols {
				g.Winner = 3 // Draw
				g.GameOver = true
			} else {
				// Switch turn
				if g.Turn == models.Red {
					g.Turn = models.Yellow
				} else {
					g.Turn = models.Red
				}
			}
			return true
		}
	}
	return false
}

// Basic Heuristic Bot
func (g *GameLogic) GetBotMove() int {
	// 1. Try to win
	for c := 0; c < models.Cols; c++ {
		if g.simulateMove(c, models.Yellow) { return c }
	}
	// 2. Block opponent
	for c := 0; c < models.Cols; c++ {
		if g.simulateMove(c, models.Red) { return c }
	}
	// 3. Random valid move
	validCols := []int{}
	for c := 0; c < models.Cols; c++ {
		if g.Board[0][c] == models.Empty { validCols = append(validCols, c) }
	}
	if len(validCols) > 0 {
		return validCols[rand.Intn(len(validCols))]
	}
	return -1
}

func (g *GameLogic) simulateMove(col, player int) bool {
	for r := models.Rows - 1; r >= 0; r-- {
		if g.Board[r][col] == models.Empty {
			g.Board[r][col] = player
			win := g.checkWin(r, col, player)
			g.Board[r][col] = models.Empty // Undo
			return win
		}
	}
	return false
}

func (g *GameLogic) checkWin(row, col, player int) bool {
	dirs := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, d := range dirs {
		count := 1
		for i := 1; i < 4; i++ {
			r, c := row+d[0]*i, col+d[1]*i
			if r < 0 || r >= models.Rows || c < 0 || c >= models.Cols || g.Board[r][c] != player { break }
			count++
		}
		for i := 1; i < 4; i++ {
			r, c := row-d[0]*i, col-d[1]*i
			if r < 0 || r >= models.Rows || c < 0 || c >= models.Cols || g.Board[r][c] != player { break }
			count++
		}
		if count >= 4 { return true }
	}
	return false
}