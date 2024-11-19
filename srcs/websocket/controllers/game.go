package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"
	"websocket/models"

	"github.com/google/uuid"
)

type GameEvent struct {
	models.Event
	LobbyId    uuid.UUID `json:"lobbyId"`
	UserId     uint64    `json:"userId"`
	State      GameState `json:"state"`
	KeyPressed string    `json:"keyPressed"`
}

type Ball struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	DX, DY float64 `json:"-"` //vitesse de la ball et direction
	Radius float64 `json:"-"`
}

type Paddle struct {
	Width            float64 `json:"width"`
	Height           float64 `json:"height"`
	Speed            float64 `json:"-"`
	Player1Y         float64 `json:"player1Y"`
	Player2Y         float64 `json:"player2Y"`
	Player1Direction int     `json:"player1YDirection"`
	Player2Direction int     `json:"player2YDirection"`
}

type Score struct {
	Player1 int `json:"player1"`
	Player2 int `json:"player2"`
}

type GameState struct {
	Ball      Ball       `json:"ball"`
	Paddles   Paddle     `json:"paddle"`
	Score     Score      `json:"score"`
	IsActive  bool       `json:"isActive"`
	Winner    uint64     `json:"winner"`
	mutex     sync.Mutex `json:"-"`
	IsPaused  bool       `json:"isPaused"`
	PauseTime time.Time  `json:"pauseTime"`
}

type Player struct {
	ID       uint64  `json:"id"`
	Position float64 `json:"position"`
}

type Game struct {
	Player1 Player    `json:"player1"`
	Player2 Player    `json:"player2"`
	State   GameState `json:"state"`
	Status  string    `json:"status"`
}

type GameCommand struct {
	PlayerID uint64
	Command  string
}

const (
	CanvasWidth         = 800
	CanvasHeight        = 600
	BallSpeed           = 5
	PaddleSpeed         = 7
	Paddle1DistanceWall = 30
	Paddle2DistanceWall = 770
	WinningScore        = 5
	paddleSpeed         = 8.0
)

// create instance of game and init all data
func NewGame(player1ID, player2ID uint64) *Game {
	return &Game{
		Player1: Player{
			ID:       player1ID,
			Position: CanvasHeight / 2,
		},

		Player2: Player{
			ID:       player2ID,
			Position: CanvasHeight / 2,
		},

		State: GameState{
			Ball: Ball{
				X:      CanvasWidth / 2,
				Y:      CanvasHeight / 2,
				DX:     BallSpeed,
				DY:     0,
				Radius: 10,
			},

			Paddles: Paddle{
				Width:    60,
				Height:   120,
				Speed:    PaddleSpeed,
				Player1Y: CanvasHeight / 2,
				Player2Y: CanvasHeight / 2,
			},

			Score: Score{
				Player1: 0,
				Player2: 0,
			},

			IsActive: true,
		},
		Status: "PREGAME",
	}
}

func (g *Game) Update() {
	g.State.mutex.Lock()
	defer g.State.mutex.Unlock()

	if !g.State.IsActive {
		return
	}

	if g.State.IsPaused {
		if time.Since(g.State.PauseTime) >= PointPauseTime {
			g.State.IsPaused = false
		} else {
			return
		}
	}
	// Update paddles
	if g.State.Paddles.Player1Direction != 0 {
		newY := g.State.Paddles.Player1Y + float64(g.State.Paddles.Player1Direction)*paddleSpeed
		g.State.Paddles.Player1Y = math.Max(0, math.Min(CanvasHeight-g.State.Paddles.Height, newY))
	}

	if g.State.Paddles.Player2Direction != 0 {
		newY := g.State.Paddles.Player2Y + float64(g.State.Paddles.Player2Direction)*paddleSpeed
		g.State.Paddles.Player2Y = math.Max(0, math.Min(CanvasHeight-g.State.Paddles.Height, newY))
	}

	// Update ball position
	g.State.Ball.X += g.State.Ball.DX
	g.State.Ball.Y += g.State.Ball.DY

	// Ball collision with top and bottom walls
	if g.State.Ball.Y <= 0 || g.State.Ball.Y >= CanvasHeight {
		g.State.Ball.DY = -g.State.Ball.DY
	}

	// Ball collision with paddles
	if g.State.Ball.X <= Paddle1DistanceWall {
		if g.State.Ball.Y >= g.State.Paddles.Player1Y &&
			g.State.Ball.Y <= g.State.Paddles.Player1Y+g.State.Paddles.Height {
			g.State.Ball.DX = BallSpeed
			g.State.Ball.DY = calculedeviation(
				g.State.Ball.Y,
				g.State.Paddles.Player1Y,
				g.State.Paddles.Height,
			)
		}
	}

	if g.State.Ball.X >= Paddle2DistanceWall {
		if g.State.Ball.Y >= g.State.Paddles.Player2Y &&
			g.State.Ball.Y <= g.State.Paddles.Player2Y+g.State.Paddles.Height {
			g.State.Ball.DX = -BallSpeed
			g.State.Ball.DY = calculedeviation(
				g.State.Ball.Y,
				g.State.Paddles.Player2Y,
				g.State.Paddles.Height,
			)
		}
	}

	// Score points
	if g.State.Ball.X <= 0 {
		g.State.Score.Player2++
		g.resetBall()
	}

	if g.State.Ball.X >= CanvasWidth {
		g.State.Score.Player1++
		g.resetBall()
	}

	// Check for winner
	if g.State.Score.Player1 == WinningScore {
		g.State.IsActive = false
		g.State.Winner = g.Player1.ID
	}

	if g.State.Score.Player2 == WinningScore {
		g.State.IsActive = false
		g.State.Winner = g.Player2.ID
	}
}

func (g *Game) resetBall() {
	g.State.Ball.X = CanvasWidth / 2
	g.State.Ball.Y = CanvasHeight / 2
	g.State.Ball.DY = 0
	g.State.IsPaused = true
	g.State.PauseTime = time.Now()

	// pour mettre la direction de balle a droite ou a gauche selon l ancien but marqué
	if g.State.Ball.DX > 0 {
		g.State.Ball.DX = -g.State.Ball.DX
	} else {
		g.State.Ball.DX = BallSpeed
	}
}

func calculedeviation(ballY, paddleY, paddleHeight float64) float64 {
	// trouve la position du milieu du paddel
	midPaddle := paddleY + (paddleHeight / 2)
	// calcule la distance du milieu du paddel a la balle
	middleDistance := midPaddle - ballY
	// calcule l angle qui est entre 1 et -1
	bounceAngle := middleDistance / (paddleHeight / 2)
	//calcule la vitesse vertival final
	verticalSpeed := -bounceAngle * (BallSpeed / 2)
	return verticalSpeed
}

func (g *Game) HandleCommand(cmd GameCommand) {
	g.State.mutex.Lock()
	defer g.State.mutex.Unlock()

	switch cmd.Command {
	case "UP":
		if cmd.PlayerID == g.Player1.ID {
			g.State.Paddles.Player1Direction = -1
		} else if cmd.PlayerID == g.Player2.ID {
			g.State.Paddles.Player2Direction = -1
		}
	case "DOWN":
		if cmd.PlayerID == g.Player1.ID {
			g.State.Paddles.Player1Direction = 1
		} else if cmd.PlayerID == g.Player2.ID {
			g.State.Paddles.Player2Direction = 1
		}
	case "STOP":
		if cmd.PlayerID == g.Player1.ID {
			g.State.Paddles.Player1Direction = 0
		} else if cmd.PlayerID == g.Player2.ID {
			g.State.Paddles.Player2Direction = 0
		}
	}
}

func handleGameMessage(h *Hub, data []byte) {
	var evt GameEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		fmt.Printf("Error GameEvent type unmarshall\n")
		return
	}
	lobby := h.Lobbies[evt.LobbyId]
	if lobby == nil {
		return
	}

	cmd := GameCommand{
		PlayerID: evt.UserId,
		Command:  evt.KeyPressed,
	}

	lobby.Game.HandleCommand(cmd)
}
