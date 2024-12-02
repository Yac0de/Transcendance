package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"websocket/models"

	"github.com/google/uuid"
)

func CreateTournamentTreeEvent(tournament *Tournament) *TournamentTreeEvent {
	return &TournamentTreeEvent{
		Event: models.Event{
			Type: "TOURNAMENT_TREE_STATE",
		},
		Code:  tournament.Id,
		Semi1: tournament.Semi1,
		Semi2: tournament.Semi2,
		Final: tournament.Final,
	}
}

func CreateGameStartEvent(tournament *Tournament, lobbyId uuid.UUID) *GameStart {
	return &GameStart{
		TournamentEvent: TournamentEvent{
			Event: models.Event{
				Type: "TOURNAMENT_GAME",
			},
			Code: tournament.Id,
		},
		LobbyId: lobbyId,
	}
}

func ShuffleTournamentOpposition(h *Hub, tournament *Tournament) {
	players := [4]uint64{tournament.Player1.Id, tournament.Player2.Id, tournament.Player3.Id, tournament.Player4.Id}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(4, func(i int, j int) { players[i], players[j] = players[j], players[i] })
	tournament.Semi1.Player1 = players[0]
	tournament.Semi1.Player2 = players[1]
	tournament.Semi2.Player1 = players[2]
	tournament.Semi2.Player2 = players[3]
}

func AppendClientToTournament(h *Hub, tournament *Tournament, request TournamentEvent) bool {
	clientJoined := h.Clients[request.UserId]
	if clientJoined == nil {
		return false
	}

	if tournament.Player2 == nil || tournament.Player2 == clientJoined {
		tournament.Player2 = clientJoined
	} else if tournament.Player3 == nil || tournament.Player3 == clientJoined {
		tournament.Player3 = clientJoined
	} else if tournament.Player4 == nil || tournament.Player4 == clientJoined {
		tournament.Player4 = clientJoined
	} else {
		SendTournamentError(h, clientJoined, request.Code, fmt.Sprintf("Tournament with code <%s> already full", request.Code))
		return false
	}

	success, _ := json.Marshal(&request)
	safeSend(clientJoined.Send, success)

	return true
}

func GetPlayerId(player *Client) uint64 {
	if player != nil {
		return player.Id
	}
	return 0
}

func GetTournament(h *Hub, code string) *Tournament {
	var tournament *Tournament = nil
	for id, tn := range h.Tournaments {
		uid := strings.Split(id, "-")[0]
		if code == uid {
			tournament = tn
		}
	}
	return tournament
}

func RefreshTournamentEvent(event *TournamentEvent, tournament *Tournament) {
	event.Player1 = GetPlayerId(tournament.Player1)
	event.Player2 = GetPlayerId(tournament.Player2)
	event.Player3 = GetPlayerId(tournament.Player3)
	event.Player4 = GetPlayerId(tournament.Player4)
}

func ClientIsPresentOnTournament(tn *Tournament, c *Client) bool {
	if c == tn.Player1 || c == tn.Player2 || c == tn.Player3 || c == tn.Player4 {
		return true
	}
	return false
}
