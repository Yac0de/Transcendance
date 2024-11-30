package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"websocket/models"

	"github.com/google/uuid"
)

type Tournament struct {
	Id          string         `json:"id"`
	Player1     *Client        `json:"player1"`
	Player2     *Client        `json:"player2"`
	Player3     *Client        `json:"player3"`
	Player4     *Client        `json:"player4"`
	Semi1       TournamentGame `json:"semi1"`
	Semi2       TournamentGame `json:"semi2"`
	Final       TournamentGame `json:"final"`
	LobbiesSemi [2]*Lobby      `json:"-"`
	LobbyFinal  *Lobby         `json:"-"`
	State       string         `json:"-"`
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

type TournamentGame struct {
	Player1    uint64   `json:"player1id"`
	Player2    uint64   `json:"player2id"`
	Score      [2]uint8 `json:"score"`
	IsFinished bool     `json:"isFinished"`
}

type TournamentTreeEvent struct {
	models.Event
	Code   string         `json:"code"`
	UserId uint64         `json:"userId"`
	Semi1  TournamentGame `json:"semi1"`
	Semi2  TournamentGame `json:"semi2"`
	Final  TournamentGame `json:"final"`
}

type TournamentTimerEvent struct {
	models.Event
	Code          string `json:"code"`
	RemainingTime int16  `json:"remainingTime"`
}

type GameStart struct {
	TournamentEvent
	LobbyId uuid.UUID `json:"lobbyId"`
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
	case "TOURNAMENT_LEAVE":
		LeaveTournament(h, request)
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
		Semi1: TournamentGame{
			Player1:    0,
			Player2:    0,
			Score:      [2]uint8{0, 0},
			IsFinished: false,
		},
		Semi2: TournamentGame{
			Player1:    0,
			Player2:    0,
			Score:      [2]uint8{0, 0},
			IsFinished: false,
		},
		Final: TournamentGame{
			Player1:    0,
			Player2:    0,
			Score:      [2]uint8{0, 0},
			IsFinished: false,
		},

		LobbiesSemi: [2]*Lobby{nil, nil},
		LobbyFinal:  nil,
		State:       "TOURNAMENT_LOBBY",
	}
}

func SendTournamentError(h *Hub, client *Client, code string, errorMessage string) {
	error := TournamentErrorEvent{
		Event: models.Event{
			Type: "TOURNAMENT_ERROR",
		},
		Code:  code,
		Error: errorMessage,
	}
	errorToBytes, _ := json.Marshal(&error)
	safeSend(client.Send, errorToBytes)
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
	tournament := GetTournament(h, request.Code)
	if tournament == nil {
		SendTournamentError(h, h.Clients[request.UserId], request.Code, fmt.Sprintf("Tournament with code <%s> does not exist", request.Code))
		return
	}

	if AppendClientToTournament(h, tournament, request) == false {
		return
	}

	RefreshTournamentEvent(&request, tournament)

	request.Code = tournament.Id
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

func LeaveWaitingLobby(h *Hub, tournament *Tournament, clientLeft *Client, request TournamentEvent) {
	if tournament.Player1 == clientLeft {
		request.Type = "TOURNAMENT_TERMINATE"
		tnTerminate, err := json.Marshal(&request)
		if err != nil {
			fmt.Printf("Impossible to parse TournamentEvent type: ", err.Error())
			return
		}
		SendDataToPlayers(tournament, tnTerminate)
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
	RefreshTournamentEvent(&request, tournament)
	jsonData, _ := json.Marshal(&request)
	SendDataToPlayers(tournament, jsonData)
}

func StartTournament(h *Hub, request TournamentEvent) {
	tournament := h.Tournaments[request.Code]
	if request.UserId != tournament.Player1.Id {
		return
	}

	if tournament.Player1 == nil || tournament.Player2 == nil || tournament.Player3 == nil || tournament.Player4 == nil {
		SendTournamentError(h, tournament.Player1, tournament.Id, "Tournament isn't full")
		return
	}

	RefreshTournamentEvent(&request, tournament)

	jsonData, _ := json.Marshal(&request)
	SendDataToPlayers(tournament, jsonData)

	go func() {
		time.Sleep(10 * time.Millisecond)
		TournamentMonitoring(h, tournament)
	}()
}

func CreateLobbyGameTournament(player1 *Client, player2 *Client) *Lobby {
	return &Lobby{
		Id:               uuid.New(),
		Sender:           player1,
		Receiver:         player2,
		Timestamps:       LobbyTimestamps{},
		PlayersReady:     [2]bool{true, true},
		IsTournamentGame: true,
	}
}

func CreateLobbies(h *Hub, tournament *Tournament) {
	ShuffleTournamentOpposition(h, tournament)

	tournament.LobbiesSemi[0] = CreateLobbyGameTournament(h.Clients[tournament.Semi1.Player1], h.Clients[tournament.Semi1.Player2])
	h.Lobbies[tournament.LobbiesSemi[0].Id] = tournament.LobbiesSemi[0]

	tournament.LobbiesSemi[1] = CreateLobbyGameTournament(h.Clients[tournament.Semi2.Player1], h.Clients[tournament.Semi2.Player2])
	h.Lobbies[tournament.LobbiesSemi[1].Id] = tournament.LobbiesSemi[1]
}

func HandleTimerEvent(h *Hub, tournament *Tournament, sec *int16) {
	event := TournamentTimerEvent{
		Event: models.Event{
			Type: "TOURNAMENT_TIMER",
		},
		Code:          tournament.Id,
		RemainingTime: *sec,
	}
	evJson, _ := json.Marshal(&event)
	SendDatasToGame(h, tournament.Semi1, evJson)
	SendDatasToGame(h, tournament.Semi2, evJson)
	*sec -= 1
	if *sec < 0 {
		if tournament.State == "TIMER_FINAL" {
			tournament.LobbyFinal = CreateLobbyGameTournament(h.Clients[tournament.Final.Player1], h.Clients[tournament.Final.Player2])
			h.Lobbies[tournament.LobbyFinal.Id] = tournament.LobbyFinal
			tournament.State = "TOURNAMENT_START_FINAL"
			return
		}
		tournament.State = "TOURNAMENT_START_SEMI"
	}
}

func PreventPlayersGameStart(tournament *Tournament, lobby *Lobby) {
	event := CreateGameStartEvent(tournament, lobby.Id)
	ev, _ := json.Marshal(&event)
	if lobby.Sender != nil {
		safeSend(lobby.Sender.Send, ev)
	} else {
		fmt.Println("PreventPlayersGameStart: Sender ", string(ev))
	}

	if lobby.Receiver != nil {
		safeSend(lobby.Receiver.Send, ev)
	} else {
		fmt.Println("PreventPlayersGameStart: Receiver ", string(ev))
	}
}

func StartSemiFinals(h *Hub, tournament *Tournament) {
	if tournament.Semi1.IsFinished == false {
		PreventPlayersGameStart(tournament, tournament.LobbiesSemi[0])
	} else {
		tree := CreateTournamentTreeEvent(tournament)
		ev, _ := json.Marshal(&tree)
		fmt.Println("SEMI 1 Event startSemiFinals: ", string(ev))
		SendDatasToGame(h, tournament.Semi1, ev)
	}
	if tournament.Semi2.IsFinished == false {
		PreventPlayersGameStart(tournament, tournament.LobbiesSemi[1])
	} else {
		tree := CreateTournamentTreeEvent(tournament)
		ev, _ := json.Marshal(&tree)
		fmt.Println("SEMI 2 Event startSemiFinals: ", string(ev))
		SendDatasToGame(h, tournament.Semi2, ev)
	}

	tournament.State = "TOURNAMENT_ON_SEMI"
	go func() {
		time.Sleep(300 * time.Millisecond)
		if !tournament.Semi1.IsFinished {
			StartRoutine(h, tournament.LobbiesSemi[0])
		}
		if !tournament.Semi2.IsFinished {
			StartRoutine(h, tournament.LobbiesSemi[1])
		}
	}()
}

func StartFinal(h *Hub, tournament *Tournament) {
	PreventPlayersGameStart(tournament, tournament.LobbyFinal)
	tournament.State = "TOURNAMENT_ON_FINAL"
	go func() {
		time.Sleep(300 * time.Millisecond)
		StartRoutine(h, tournament.LobbyFinal)
	}()
}

func TournamentMonitoring(h *Hub, tournament *Tournament) {
	gameTicker := time.NewTicker(time.Second)
	tournament.State = "TIMER_SEMI_FINAL"
	sec := int16(15)

	CreateLobbies(h, tournament)
	event := CreateTournamentTreeEvent(tournament)
	jsonData, _ := json.Marshal(&event)
	SendDataToPlayers(tournament, jsonData)
	go func() {
		for {
			select {
			case <-gameTicker.C:
				if strings.HasPrefix(tournament.State, "TIMER_") && sec >= 0 {
					HandleTimerEvent(h, tournament, &sec)
				} else if tournament.State == "TOURNAMENT_START_SEMI" {
					StartSemiFinals(h, tournament)
				} else if tournament.State == "TOURNAMENT_ON_SEMI" {
					UpdateSemiFinals(h, tournament, event)
					if tournament.Final.Player1 != 0 && tournament.Final.Player2 != 0 {
						sec = 15
						tournament.State = "TIMER_FINAL"
					}
				} else if tournament.State == "TOURNAMENT_START_FINAL" {
					StartFinal(h, tournament)
				} else if IsFinalTournamentFinished(tournament) {
					SendTournamentTreeState(h, tournament, event)
					delete(h.Tournaments, tournament.Id)
					return
				}
			}
		}
	}()
}

func IsFinalTournamentFinished(tournament *Tournament) bool {
	return tournament.State == "TOURNAMENT_ON_FINAL" && tournament.LobbyFinal.Game.State.Winner != 0
}

func SendTournamentTreeState(h *Hub, tournament *Tournament, event *TournamentTreeEvent) {
	event.Final.Score[0] = uint8(tournament.LobbyFinal.Game.State.Score.Player1)
	event.Final.Score[1] = uint8(tournament.LobbyFinal.Game.State.Score.Player2)
	event.Final.IsFinished = true
	jsonData, _ := json.Marshal(&event)
	SendDatasToGame(h, tournament.Semi1, jsonData)
	SendDatasToGame(h, tournament.Semi2, jsonData)
}

func SendDatasToGame(h *Hub, players TournamentGame, message []byte) {
	player1 := h.Clients[players.Player1]
	player2 := h.Clients[players.Player2]
	if player1 != nil {
		safeSend(player1.Send, message)
	}

	if player2 != nil {
		safeSend(player2.Send, message)
	}
}

func UpdateGameTreeEvent(tn *Tournament, evt *TournamentTreeEvent) {
	evt.Semi1.Player1 = tn.Semi1.Player1
	evt.Semi1.Player2 = tn.Semi1.Player2
	evt.Semi1.Score = tn.Semi1.Score
	evt.Semi1.IsFinished = tn.Semi1.IsFinished
	evt.Semi2.Player1 = tn.Semi2.Player1
	evt.Semi2.Player2 = tn.Semi2.Player2
	evt.Semi2.Score = tn.Semi2.Score
	evt.Semi2.IsFinished = tn.Semi2.IsFinished
	evt.Final.Player1 = tn.Final.Player1
	evt.Final.Player2 = tn.Final.Player2
	evt.Final.Score = tn.Final.Score
	evt.Final.IsFinished = tn.Final.IsFinished
}

func UpdateSemiFinals(h *Hub, tournament *Tournament, event *TournamentTreeEvent) {
	if tournament.Final.Player1 == 0 && tournament.LobbiesSemi[0].Game.State.Winner != 0 {
		state := tournament.LobbiesSemi[0].Game.State
		tournament.Final.Player1 = state.Winner
		tournament.Semi1.Score = [2]uint8{state.Score.Player1, state.Score.Player2}
		tournament.Semi1.IsFinished = true
		UpdateGameTreeEvent(tournament, event)
		UpdateGameTreeEvent(tournament, event)
		jsonData, _ := json.Marshal(&event)
		SendDatasToGame(h, tournament.Semi1, jsonData)
		SendDatasToGame(h, tournament.Semi2, jsonData)

	}
	if tournament.Final.Player2 == 0 && tournament.LobbiesSemi[1].Game.State.Winner != 0 {
		state := tournament.LobbiesSemi[1].Game.State
		tournament.Semi2.Score = [2]uint8{state.Score.Player1, state.Score.Player2}
		tournament.Semi2.IsFinished = true
		tournament.Final.Player2 = state.Winner
		UpdateGameTreeEvent(tournament, event)
		UpdateGameTreeEvent(tournament, event)
		jsonData, _ := json.Marshal(&event)
		SendDatasToGame(h, tournament.Semi1, jsonData)
		SendDatasToGame(h, tournament.Semi2, jsonData)
	}
}

func TournamentClientHasLeft(h *Hub, tn *Tournament, c *Client) {
	evt := TournamentEvent{
		Event: models.Event{
			Type: "TOURNAMENT_LEAVE",
		},
		Code:   tn.Id,
		UserId: c.Id,
	}
	LeaveTournament(h, evt)
	return
}

func LeaveTournament(h *Hub, request TournamentEvent) {
	clientLeft := h.Clients[request.UserId]
	tournament := h.Tournaments[request.Code]
	if clientLeft == nil || tournament == nil {
		return
	}

	if tournament.State == "TOURNAMENT_LOBBY" {
		LeaveWaitingLobby(h, tournament, clientLeft, request)
	} else if tournament.State == "TOURNAMENT_ON_SEMI" {
		if tournament.Semi1.Player1 == clientLeft.Id || tournament.Semi1.Player2 == clientLeft.Id {
			tournament.LobbiesSemi[0].Game.PlayerLeaved(clientLeft.Id)
		} else if tournament.Semi2.Player1 == clientLeft.Id || tournament.Semi2.Player2 == clientLeft.Id {
			tournament.LobbiesSemi[1].Game.PlayerLeaved(clientLeft.Id)
		}
	} else if tournament.State == "TOURNAMENT_ON_FINAL" {
		if tournament.Final.Player1 == clientLeft.Id || tournament.Final.Player2 == clientLeft.Id {
			tournament.LobbyFinal.Game.PlayerLeaved(clientLeft.Id)
		}
	} else if strings.HasPrefix(tournament.State, "TIMER_SEMI_FINAL") {
		if tournament.Semi1.Player1 == clientLeft.Id {
			tournament.Semi1.Score = [2]uint8{0, 1}
			tournament.Semi1.IsFinished = true
			tournament.Final.Player1 = tournament.Semi1.Player2
			fmt.Println("Game Semi1: ", tournament)
		} else if tournament.Semi1.Player2 == clientLeft.Id {
			tournament.Semi1.Score = [2]uint8{1, 0}
			tournament.Semi1.IsFinished = true
			tournament.Final.Player1 = tournament.Semi1.Player1
			fmt.Println("Game Semi1: ", tournament)
		} else if tournament.Semi2.Player1 == clientLeft.Id {
			tournament.Semi2.Score = [2]uint8{0, 1}
			tournament.Semi2.IsFinished = true
			tournament.Final.Player2 = tournament.Semi2.Player2
			fmt.Println("Game Semi2: ", tournament)
		} else if tournament.Semi2.Player2 == clientLeft.Id {
			tournament.Semi2.Score = [2]uint8{1, 0}
			tournament.Semi2.IsFinished = true
			tournament.Final.Player2 = tournament.Semi2.Player1
			fmt.Println("Game Semi2: ", tournament)
		}
	}
}
