package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
	"websocket/models"

	"github.com/google/uuid"
)

type Tournament struct {
	Id      string        `json:"id"`
	Player1 *Client       `json:"player1"`
	Player2 *Client       `json:"player2"`
	Player3 *Client       `json:"player3"`
	Player4 *Client       `json:"player4"`
	Game1   [2]uint64     `json:"game1"`
	Game2   [2]uint64     `json:"game2"`
	Mutex   sync.Mutex    `json:"-"`
	Destroy chan struct{} `json:"-"`
}

type TournamentEvent struct {
	models.Event
	Code    string    `json:"code"`
	UserId  uint64    `json:"userId"`
	Player1 uint64    `json:"player1id"`
	Player2 uint64    `json:"player2id"`
	Player3 uint64    `json:"player3id"`
	Player4 uint64    `json:"player4id"`
	Game1   [2]uint64 `json:"game1"`
	Game2   [2]uint64 `json:"game2"`
}

type TournamentTimerEvent struct {
	models.Event
	Code          string `json:"code"`
	RemainingTime int16  `json:"remainingTime"`
}

type TournamentErrorEvent struct {
	models.Event
	Code  string `json:"code"`
	Error string `json:"error"`
}

func HandleTournament(h *Hub, event string, data []byte) {
	var request TournamentEvent
	if err := json.Unmarshal(data, &request); err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}
	switch event {
	case "TOURNAMENT_CREATE":
		CreateTournament(h, request)
	case "TOURNAMENT_JOIN_WITH_CODE":
		JoinTournament(h, request)
	case "TOURNAMENT_LEAVE_WAITING_ROOM":
		LeaveWaitingRoomTournament(h, request)
	case "TOURNAMENT_START":
		StartTournament(h, request)
	}
}

func NewTournament(h *Hub, request TournamentEvent) *Tournament {
	return &Tournament{
		Id:      uuid.New().String(),
		Player1: h.Clients[request.UserId],
		Player2: nil,
		Player3: nil,
		Player4: nil,
		Game1:   [2]uint64{0, 0},
		Game2:   [2]uint64{0, 0},
	}
}

func SendTournamentError(h *Hub, request TournamentEvent, errorMessage string) {
	error := TournamentErrorEvent{
		Event: models.Event{
			Type: "TOURNAMENT_ERROR",
		},
		Code:  request.Code,
		Error: errorMessage,
	}
	errorToBytes, _ := json.Marshal(&error)
	tournament, exist := h.Tournaments[request.Code]
	if !exist {
		safeSend(h.Clients[request.UserId].Send, errorToBytes)
		return
	}
	SendDataToPlayers(tournament, errorToBytes)
}

func CreateTournament(h *Hub, request TournamentEvent) {
	tournament := NewTournament(h, request)
	fmt.Println(tournament)
	h.Tournaments[tournament.Id] = tournament
	request.Player1 = tournament.Player1.Id
	request.Code = tournament.Id
	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}
	safeSend(tournament.Player1.Send, jsonData)
}

func JoinTournament(h *Hub, request TournamentEvent) {
	var tournament *Tournament = nil
	for id, tn := range h.Tournaments {
		uid := strings.Split(id, "-")[0]
		if request.Code == uid {
			tournament = tn
		}
	}
	if tournament == nil {
		SendTournamentError(h, request, fmt.Sprintf("Tournament with code <%s> does not exist", request.UserId))
		return
	}
	clientJoined := h.Clients[request.UserId]
	if tournament.Player2 == nil || tournament.Player2 == clientJoined {
		tournament.Player2 = clientJoined
	} else if tournament.Player3 == nil || tournament.Player3 == clientJoined {
		tournament.Player3 = clientJoined
	} else if tournament.Player4 == nil || tournament.Player4 == clientJoined {
		tournament.Player4 = clientJoined
	} else {
		SendTournamentError(h, request, fmt.Sprintf("Tournament with code <%s> already full", request.Code))
		return
	}
	success, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}
	safeSend(clientJoined.Send, success)

	request.Code = tournament.Id

	request.Player1 = getPlayerId(tournament.Player1)
	request.Player2 = getPlayerId(tournament.Player2)
	request.Player3 = getPlayerId(tournament.Player3)
	request.Player4 = getPlayerId(tournament.Player4)

	request.Type = "TOURNAMENT_EVENT"
	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}

	go func() {
		time.Sleep(10 * time.Millisecond)
		SendDataToPlayers(tournament, jsonData)
	}()
}

func getPlayerId(player *Client) uint64 {
	if player != nil {
		return player.Id
	}
	return 0
}

func SendDataToPlayers(tournament *Tournament, datas []byte) {
	if tournament.Player1 != nil {
		safeSend(tournament.Player1.Send, datas)
	}
	if tournament.Player2 != nil {
		safeSend(tournament.Player2.Send, datas)
	}
	if tournament.Player3 != nil {
		safeSend(tournament.Player3.Send, datas)
	}
	if tournament.Player4 != nil {
		safeSend(tournament.Player4.Send, datas)
	}
}

func LeaveWaitingRoomTournament(h *Hub, request TournamentEvent) {
	clientLeft := h.Clients[request.UserId]
	tournament := h.Tournaments[request.Code]
	if clientLeft == nil || tournament == nil {
		return
	}
	if tournament.Player1 == clientLeft {
		request.Type = "TOURNAMENT_TERMINATE"
		jsonData, err := json.Marshal(&request)
		if err != nil {
			fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
			return
		}
		SendDataToPlayers(tournament, jsonData)
		delete(h.Tournaments, tournament.Id)
		return
	} else if tournament.Player2 == clientLeft {
		tournament.Player2 = nil
	} else if tournament.Player3 == clientLeft {
		tournament.Player3 = nil
	} else if tournament.Player4 == clientLeft {
		tournament.Player4 = nil
	}

	request.Type = "TOURNAMENT_EVENT"
	request.Player1 = getPlayerId(tournament.Player1)
	request.Player2 = getPlayerId(tournament.Player2)
	request.Player3 = getPlayerId(tournament.Player3)
	request.Player4 = getPlayerId(tournament.Player4)

	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}
	SendDataToPlayers(tournament, jsonData)
}

func StartTournament(h *Hub, request TournamentEvent) {
	tournament := h.Tournaments[request.Code]
	if request.UserId != tournament.Player1.Id {
		SendTournamentError(h, request, "Only the creator can start the tournament")
		return
	}

	// if tournament.Player1 == nil || tournament.Player2 == nil || tournament.Player3 == nil || tournament.Player4 == nil {
	// 	SendTournamentError(h, request, "Tournament is not full")
	// 	return
	// }

	RefreshTournamentEvent(&request, tournament)

	tournament.Game1[0] = request.Player1
	tournament.Game1[1] = request.Player2
	tournament.Game2[0] = request.Player3
	tournament.Game2[1] = request.Player4

	request.Game1 = tournament.Game1
	request.Game2 = tournament.Game2
	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}

	SendDataToPlayers(tournament, jsonData)

	go func() {
		time.Sleep(10 * time.Millisecond)
		TournamentMonitoring(h, tournament)
	}()
	return
}

func TournamentMonitoring(h *Hub, tournament *Tournament) {
	gameTicker := time.NewTicker(time.Second)
	state := "TIMER"
	sec := int16(15)
	go func() {
		for {
			select {
			case <-tournament.Destroy:
				gameTicker.Stop()
				return
			case <-gameTicker.C:
				if state == "TIMER" && sec >= 0 {
					event := TournamentTimerEvent{
						Event: models.Event{
							Type: "TOURNAMENT_TIMER",
						},
						Code:          tournament.Id,
						RemainingTime: sec,
					}
					evJson, _ := json.Marshal(&event)
					SendDataToPlayers(tournament, evJson)
					sec -= 1
					if sec < 0 {
						state = "GAME_START"
					}
				} else if state == "GAME_START" {
					event := TournamentTimerEvent{
						Event: models.Event{
							Type: "TOURNAMENT_GAME",
						},
						Code: tournament.Id,
					}
					evJson, _ := json.Marshal(&event)
					SendDataToPlayers(tournament, evJson)
					state = "TOURNAMENT_ON_GAME"
				}
			}
		}
	}()

}

func RefreshTournamentEvent(event *TournamentEvent, tournament *Tournament) {
	event.Player1 = getPlayerId(tournament.Player1)
	event.Player2 = getPlayerId(tournament.Player2)
	event.Player3 = getPlayerId(tournament.Player3)
	event.Player4 = getPlayerId(tournament.Player4)
}

func ClientIsPresentOnTournament(tn *Tournament, c *Client) bool {
	if c == tn.Player1 || c == tn.Player2 || c == tn.Player3 || c == tn.Player4 {
		return true
	}
	return false
}
func TournamentClientHasLeft(h *Hub, tn *Tournament, c *Client) {
	eventName := "TOURNAMENT_EVENT"
	if c == tn.Player1 {
		eventName = "TOURNAMENT_TERMINATE"
		tn.Player1 = nil
	} else if c == tn.Player2 {
		tn.Player2 = nil
	} else if c == tn.Player3 {
		tn.Player3 = nil
	} else if c == tn.Player4 {
		tn.Player4 = nil
	}

	event := TournamentEvent{
		Event: models.Event{
			Type: eventName,
		},
		Code: tn.Id,
	}

	RefreshTournamentEvent(&event, tn)

	jsonData, err := json.Marshal(&event)
	if err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}

	SendDataToPlayers(tn, jsonData)

	if c == tn.Player1 {
		delete(h.Tournaments, tn.Id)
	}
}
