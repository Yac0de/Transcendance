package controllers

import (
	"encoding/json"
	"fmt"
	"websocket/models"
)

type Tournament struct {
	Id      string `json:"id"`
	Player1 uint64 `json:"player1"`
	Player2 uint64 `json:"player2"`
	Player3 uint64 `json:"player3"`
	Player4 uint64 `json:"player4"`
	// Mutex   sync.Mutex    `json:"-"`
	// Destroy chan struct{} `json:"-"`
}

type TournamentEvent struct {
	models.Event
	Code    string `json:"code"`
	UserId  uint64 `json:"userId"`
	Player1 uint64 `json:"player1"`
	Player2 uint64 `json:"player2"`
	Player3 uint64 `json:"player3"`
	Player4 uint64 `json:"player4"`
}

func HandleTournament(h *Hub, event string, data []byte) {
	var request TournamentEvent
	if err := json.Unmarshal(data, &request); err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}
	switch event {
	case "TOURNAMENT_CREATE":
	}
}

func NewTournament() *Tournament {
	return &Tournament{
		Id:      "",
		Player1: 0,
		Player2: 0,
		Player3: 0,
		Player4: 0,
	}
}
