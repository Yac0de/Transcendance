package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"websocket/models"

	"github.com/google/uuid"
)

type Tournament struct {
	Id      string  `json:"id"`
	Player1 *Client `json:"player1"`
	Player2 *Client `json:"player2"`
	Player3 *Client `json:"player3"`
	Player4 *Client `json:"player4"`
}

type TournamentEvent struct {
	models.Event
	Code    string `json:"code"`
	UserId  uint64 `json:"userId"`
	Player1 uint64 `json:"player1id"`
	Player2 uint64 `json:"player2id"`
	Player3 uint64 `json:"player3id"`
	Player4 uint64 `json:"player4id"`
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
	return &Tournament{
		Id:      uuid.New().String(),
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
	tournament := NewTournament(h, request)
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
	success, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
		return
	}
	safeSend(clientJoined.Send, success)

	request.Player1 = tournament.Player1.Id
	request.Player2 = tournament.Player2.Id
	request.Player3 = tournament.Player3.Id
	request.Player4 = tournament.Player4.Id
	request.Type = "TOURNAMENT_EVENT"

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
