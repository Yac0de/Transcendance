package controllers

import (
	"encoding/json"
	"fmt"
	"websocket/models"
)

type Tournament struct {
	Id      string  `json:"id"`
	Player1 *Client `json:"player1"`
	Player2 *Client `json:"player2"`
	Player3 *Client `json:"player3"`
	Player4 *Client `json:"player4"`
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
	}
}

func NewTournament(h *Hub, request TournamentEvent) *Tournament {
	if _, exist := h.Tournaments[request.Code]; exist {
		return nil
	}
	return &Tournament{
		Id:      request.Code,
		Player1: h.Clients[request.UserId],
		Player2: nil,
		Player3: nil,
		Player4: nil,
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
	if _, exist := h.Clients[request.UserId]; !exist {
		SendTournamentError(h, request, fmt.Sprintf("Client id <%s> does not log", request.UserId))
		return
	}

	tournament := NewTournament(h, request)
	if tournament == nil {
		SendTournamentError(h, request, fmt.Sprintf("Tournament <%s> code already exist.", request.Code))
		return
	}
	h.Tournaments[tournament.Id] = tournament
	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}
	safeSend(tournament.Player1.Send, jsonData)
}

func JoinTournament(h *Hub, request TournamentEvent) {
	tournament, exist := h.Tournaments[request.Code]
	if !exist {
		SendTournamentError(h, request, fmt.Sprintf("Tournament with code <%s> does not exist", request.UserId))
		return
	}
	clientJoined := h.Clients[request.UserId]
	if tournament.Player2 == nil {
		tournament.Player2 = clientJoined
	} else if tournament.Player3 == nil {
		tournament.Player3 = clientJoined
	} else if tournament.Player4 == nil {
		tournament.Player4 = clientJoined
	} else {
		SendTournamentError(h, request, fmt.Sprintf("Tournament with code <%s> already full", request.Code))
		return
	}
	request.Player1 = tournament.Player1.Id
	request.Player2 = tournament.Player2.Id
	request.Player3 = tournament.Player3.Id
	request.Player4 = tournament.Player4.Id

	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}
	SendDataToPlayers(tournament, jsonData)
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
