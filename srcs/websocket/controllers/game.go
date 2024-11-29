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
	LobbyId          uuid.UUID `json:"lobbyId"`
	UserId           uint64    `json:"userId"`
	State            GameState `json:"state"`
	KeyPressed       string    `json:"keyPressed"`
	Player1Id        uint64    `json:"player1id"`
	Player2Id        uint64    `json:"player2id"`
	IsTournamentGame bool      `json:"isTournamentGame"`
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
	Player1X         float64 `json:"player1X"`
	Player2X         float64 `json:"player2X"`
	Player1Direction int     `json:"player1YDirection"`
	Player2Direction int     `json:"player2YDirection"`
}

type Score struct {
	Player1 int `json:"player1"`
	Player2 int `json:"player2"`
}

type GameState struct {
	Ball          Ball       `json:"ball"`
	Paddles       Paddle     `json:"paddle"`
	Score         Score      `json:"score"`
	IsActive      bool       `json:"isActive"`
	Winner        uint64     `json:"winner"`
	mutex         sync.Mutex `json:"-"`
	IsPaused      bool       `json:"isPaused"`
	PauseTime     time.Time  `json:"pauseTime"`
	Player1Boost  BoostState `json:"player1boost"`
	Player2Boost  BoostState `json:"player2boost"`
	RemainingTime int        `json:"remainingTime"`
}

type BoostState struct {
	BallHit       int  `json:"ballhit"`
	BoostReady    bool `json:"boostReady"`
	IsBoostActive bool `json:"isboostactive"`
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
	BallSpeed           = 8
	PaddleSpeed         = 3
	Paddle1DistanceWall = 20
	Paddle2DistanceWall = 760
	WinningScore        = 1
	paddleSpeed         = 8.0
	collisionToBoost    = 3
	boostMultiplier     = 2.5
)

// create instance of game and init all data
func NewGame(player1ID uint64, player2ID uint64) *Game {
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
				Width:    20,
				Height:   120,
				Speed:    PaddleSpeed,
				Player1Y: (CanvasHeight / 2) - 120/2,
				Player2Y: (CanvasHeight / 2) - 120/2,
				Player1X: Paddle1DistanceWall,
				Player2X: Paddle2DistanceWall,
			},

			Score: Score{
				Player1: 0,
				Player2: 0,
			},
			Player1Boost: BoostState{
				BallHit:       0,
				BoostReady:    false,
				IsBoostActive: false,
			},
			Player2Boost: BoostState{
				BallHit:       0,
				BoostReady:    false,
				IsBoostActive: false,
			},
			Winner:   0,
			IsActive: true,

			RemainingTime: 300,
		},
		Status: "PREGAME",
	}
}

func (g *Game) PlayerLeaved(id uint64) {
	g.State.mutex.Lock()
	defer g.State.mutex.Unlock()

	if !g.State.IsActive {
		return
	}

	if id == g.Player2.ID {
		g.State.Winner = g.Player1.ID
		g.State.Score.Player1 = WinningScore
		g.State.Score.Player2 = 0
		return
	}
	g.State.Winner = g.Player2.ID
	g.State.Score.Player2 = WinningScore
	g.State.Score.Player1 = 0
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

	now := time.Now()
	if now.Sub(g.State.PauseTime) >= time.Second {
		g.State.PauseTime = now
		if g.State.RemainingTime > 0 {
			g.State.RemainingTime--
		} else {
			g.State.IsActive = false
			if g.State.Score.Player1 > g.State.Score.Player2 {
				g.State.Winner = g.Player1.ID
			} else if g.State.Score.Player2 > g.State.Score.Player1 {
				g.State.Winner = g.Player2.ID
			}
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
	if g.State.Ball.Y-g.State.Ball.Radius <= 0 || g.State.Ball.Y+g.State.Ball.Radius >= CanvasHeight {
		g.State.Ball.DY = -g.State.Ball.DY
	}

	// Classic ball collision with paddle 1
	if g.State.Ball.X <= Paddle1DistanceWall+g.State.Paddles.Width+g.State.Ball.Radius {
		if g.State.Ball.Y >= g.State.Paddles.Player1Y &&
			g.State.Ball.Y <= g.State.Paddles.Player1Y+g.State.Paddles.Height &&
			g.State.Ball.X-g.State.Ball.Radius > g.State.Paddles.Player1X {

			multiplier := 1.0
			if g.State.Player1Boost.IsBoostActive {
				multiplier = boostMultiplier
				g.State.Player1Boost.IsBoostActive = false
				g.State.Player1Boost.BoostReady = false
			}

			g.State.Ball.DX = BallSpeed * multiplier
			g.State.Ball.DY = computeDeviation(
				g.State.Ball.Y,
				g.State.Paddles.Player1Y,
				g.State.Paddles.Height,
			) * multiplier
			g.hitCounter(1)
		}
	}

	// Classic ball collision with paddle 2
	if g.State.Ball.X >= Paddle2DistanceWall-g.State.Ball.Radius {
		if g.State.Ball.Y >= g.State.Paddles.Player2Y &&
			g.State.Ball.Y <= g.State.Paddles.Player2Y+g.State.Paddles.Height &&
			g.State.Ball.X+g.State.Ball.Radius < g.State.Paddles.Player2X+g.State.Paddles.Width {

			multiplier := 1.0
			if g.State.Player2Boost.IsBoostActive {
				multiplier = boostMultiplier
				g.State.Player2Boost.IsBoostActive = false
				g.State.Player2Boost.BoostReady = false
			}
			g.State.Ball.DX = -BallSpeed * multiplier // Negative because ball should go left
			g.State.Ball.DY = computeDeviation(
				g.State.Ball.Y,
				g.State.Paddles.Player2Y,
				g.State.Paddles.Height,
			) * multiplier
			g.hitCounter(2)
		}
	}
	//Top part of the paddle collision
	if g.isBallAbovePaddle() {
		if g.State.Ball.X+g.State.Ball.Radius >= g.State.Paddles.Player1X &&
			g.State.Ball.X-g.State.Ball.Radius <= g.State.Paddles.Player1X+g.State.Paddles.Width {

			// Calculate vertical distance between ball and paddle top edge
			distanceY := math.Abs(g.State.Ball.Y - g.State.Paddles.Player1Y)

			// If distance is less than ball radius, we have a collision
			if distanceY-5 <= g.State.Ball.Radius {
				overlap := g.State.Ball.Radius - distanceY
				g.State.Ball.DY = g.State.Ball.Y - overlap - 1
				g.State.Ball.DX = computeSideDeviation(
					g.State.Ball.X,
					g.State.Paddles.Player1X,
					g.State.Paddles.Width,
				)
				g.State.Ball.DY = -BallSpeed

			}
		}
	}

	if g.isBallBelowPaddle() {
		if g.State.Ball.X+g.State.Ball.Radius >= g.State.Paddles.Player1X &&
			g.State.Ball.X-g.State.Ball.Radius <= g.State.Paddles.Player1X+g.State.Paddles.Width {
			// Calculate vertical distance between ball and paddle bottom edge
			distanceY := math.Abs(g.State.Ball.Y - (g.State.Paddles.Player1Y + g.State.Paddles.Height))
			// If distance is less than ball radius, we have a collision
			if distanceY-5 <= g.State.Ball.Radius {
				overlap := g.State.Ball.Radius - distanceY
				g.State.Ball.DY = g.State.Ball.Y + overlap + 1
				g.State.Ball.DX = computeSideDeviation(
					g.State.Ball.X,
					g.State.Paddles.Player1X,
					g.State.Paddles.Width,
				)
				g.State.Ball.DY = BallSpeed

			}
		}
	}

	//Top part of the paddle collision for Player 2
	if g.isBallAbovePaddle() {
		if g.State.Ball.X+g.State.Ball.Radius >= g.State.Paddles.Player2X &&
			g.State.Ball.X-g.State.Ball.Radius <= g.State.Paddles.Player2X+g.State.Paddles.Width {
			// Calculate vertical distance between ball and paddle top edge
			distanceY := math.Abs(g.State.Ball.Y - g.State.Paddles.Player2Y)
			// If distance is less than ball radius, we have a collision
			if distanceY-5 <= g.State.Ball.Radius {
				overlap := g.State.Ball.Radius - distanceY
				g.State.Ball.DY = g.State.Ball.Y - overlap - 1
				g.State.Ball.DX = computeSideDeviation(
					g.State.Ball.X,
					g.State.Paddles.Player2X,
					g.State.Paddles.Width,
				)
				g.State.Ball.DY = -BallSpeed

			}
		}
	}
	if g.isBallBelowPaddle() {
		if g.State.Ball.X+g.State.Ball.Radius >= g.State.Paddles.Player2X &&
			g.State.Ball.X-g.State.Ball.Radius <= g.State.Paddles.Player2X+g.State.Paddles.Width {
			// Calculate vertical distance between ball and paddle bottom edge
			distanceY := math.Abs(g.State.Ball.Y - (g.State.Paddles.Player2Y + g.State.Paddles.Height))
			// If distance is less than ball radius, we have a collision
			if distanceY-5 <= g.State.Ball.Radius {
				overlap := g.State.Ball.Radius - distanceY
				g.State.Ball.DY = g.State.Ball.Y + overlap + 1
				g.State.Ball.DX = computeSideDeviation(
					g.State.Ball.X,
					g.State.Paddles.Player2X,
					g.State.Paddles.Width,
				)
				g.State.Ball.DY = BallSpeed

			}
		}
	}

	// Score points
	if g.State.Ball.X <= 0 {
		g.State.Score.Player2++
		g.resetBall()
		g.resetPaddle()
	}

	if g.State.Ball.X >= CanvasWidth {
		g.State.Score.Player1++
		g.resetBall()
		g.resetPaddle()

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

func (g *Game) isBallAbovePaddle() bool {
	if g.State.Ball.Y+g.State.Ball.Radius <= g.State.Paddles.Player1Y {
		return true
	} else {
		return false
	}
}

func (g *Game) isBallBelowPaddle() bool {
	if g.State.Ball.Y-g.State.Ball.Radius >= g.State.Paddles.Player1Y+g.State.Paddles.Height {
		return true
	} else {
		return false
	}
}

func computeDeviation(ballY, paddleY, paddleHeight float64) float64 {
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

func computeSideDeviation(ballX, paddleX, paddleWidth float64) float64 {
	// trouve la position du milieu du paddle horizontalement
	midPaddle := paddleX + (paddleWidth / 2)
	// calcule la distance du milieu du paddle à la balle horizontalement
	middleDistance := midPaddle - ballX
	// calcule l'angle qui est entre 1 et -1
	bounceAngle := middleDistance / (paddleWidth / 2)
	// calcule la vitesse horizontale finale
	horizontalSpeed := -bounceAngle * (BallSpeed / 2)
	return horizontalSpeed
}

func (g *Game) resetBall() {
	g.State.Ball.X = CanvasWidth / 2
	g.State.Ball.Y = CanvasHeight / 2
	g.State.Ball.DY = 0
	g.State.IsPaused = true
	g.State.PauseTime = time.Now()
	g.State.Player1Boost.BallHit = 0
	g.State.Player1Boost.BoostReady = false
	g.State.Player1Boost.IsBoostActive = false
	g.State.Player2Boost.BallHit = 0
	g.State.Player2Boost.BoostReady = false
	g.State.Player2Boost.IsBoostActive = false

	if g.State.Ball.DX > 0 {
		g.State.Ball.DX = -BallSpeed
	} else {
		g.State.Ball.DX = BallSpeed
	}
}

func (g *Game) resetPaddle() {
	g.State.Paddles.Player1Y = (CanvasHeight / 2) - 120/2
	g.State.Paddles.Player2Y = (CanvasHeight / 2) - 120/2
}

func (g *Game) hitCounter(playerNum int) {
	if playerNum == 1 {
		g.State.Player1Boost.BallHit++
		if g.State.Player1Boost.BallHit >= collisionToBoost {
			g.State.Player1Boost.BoostReady = true
			g.State.Player1Boost.BallHit = 0
		}
	} else {
		g.State.Player2Boost.BallHit++
		if g.State.Player2Boost.BallHit >= collisionToBoost {
			g.State.Player2Boost.BoostReady = true
			g.State.Player2Boost.BallHit = 0
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
	fmt.Printf("CMD: %+v\n", cmd)

	lobby.Game.HandleCommand(cmd)
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
	case "SPACE":
		if cmd.PlayerID == g.Player1.ID && g.State.Player1Boost.BoostReady {
			g.State.Player1Boost.BoostReady = false
			g.State.Player1Boost.IsBoostActive = true
		} else if cmd.PlayerID == g.Player2.ID && g.State.Player2Boost.BoostReady {
			g.State.Player2Boost.BoostReady = false
			g.State.Player2Boost.IsBoostActive = true
		}
	}
}
