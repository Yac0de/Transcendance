package controllers

import (
	"encoding/json"
	"fmt"
	"time"
	"websocket/models"

	"github.com/google/uuid"
)

type Tournament struct {
	Id          string    `json:"id"`
	Player1     *Client   `json:"player1"`
	Player2     *Client   `json:"player2"`
	Player3     *Client   `json:"player3"`
	Player4     *Client   `json:"player4"`
	Semi1       [2]uint64 `json:"semi1"`
	Semi2       [2]uint64 `json:"semi2"`
	Final       [2]uint64 `json:"final"`
	LobbiesSemi [2]*Lobby `json:"-"`
	LobbyFinal  *Lobby    `json:"-"`
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
	case "TOURNAMENT_LEAVE_WAITING_ROOM":
		LeaveWaitingRoomTournament(h, request)
	case "TOURNAMENT_START":
		StartTournament(h, request)
	}
}

func NewTournament(h *Hub, request TournamentEvent) *Tournament {
	return &Tournament{
		Id:          uuid.New().String(),
		Player1:     h.Clients[request.UserId],
		Player2:     nil,
		Player3:     nil,
		Player4:     nil,
		Semi1:       [2]uint64{0, 0},
		Semi2:       [2]uint64{0, 0},
		Final:       [2]uint64{0, 0},
		LobbiesSemi: [2]*Lobby{nil, nil},
		LobbyFinal:  nil,
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
	RefreshTournamentEvent(&request, tournament)

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

	semi1 := CreateLobbyGameTournament(h.Clients[tournament.Semi1[0]], h.Clients[tournament.Semi1[1]])
	semi2 := CreateLobbyGameTournament(h.Clients[tournament.Semi2[0]], h.Clients[tournament.Semi2[1]])

	h.Lobbies[semi1.Id] = semi1
	tournament.LobbiesSemi[0] = semi1

	h.Lobbies[semi2.Id] = semi2
	tournament.LobbiesSemi[1] = semi2
}

func TournamentMonitoring(h *Hub, tournament *Tournament) {
	gameTicker := time.NewTicker(time.Second)
	state := "TIMER"
	sec := int16(3)

	CreateLobbies(h, tournament)

	event := TournamentTreeEvent{
		Event: models.Event{
			Type: "TOURNAMENT_TREE_STATE",
		},
		Code: tournament.Id,
		Semi1: TournamentGame{
			Player1:    tournament.Semi1[0],
			Player2:    tournament.Semi1[1],
			Score:      [2]uint8{0, 0},
			IsFinished: false,
		},
		Semi2: TournamentGame{
			Player1:    tournament.Semi2[0],
			Player2:    tournament.Semi2[1],
			Score:      [2]uint8{0, 0},
			IsFinished: false,
		},
		Final: TournamentGame{
			Player1:    0,
			Player2:    0,
			Score:      [2]uint8{0, 0},
			IsFinished: false,
		},
	}
	jsonData, _ := json.Marshal(&event)
	SendDataToPlayers(tournament, jsonData)
	go func() {
		for {
			select {
			case <-gameTicker.C:
				if (state == "TIMER" || state == "TIMER_FINAL") && sec >= 0 {
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
						fmt.Printf("TIMER SEC FINISH\n")
						if state == "TIMER_FINAL" {
							tournament.LobbyFinal = CreateLobbyGameTournament(h.Clients[tournament.Final[0]], h.Clients[tournament.Final[1]])
							h.Lobbies[tournament.LobbyFinal.Id] = tournament.LobbyFinal
							fmt.Printf("TIMER_FINAL -> TOURNAMENT_START_FINAL: %+v\n", tournament.LobbyFinal)
							fmt.Printf("Lobbies: %+v\n", h.Lobbies)
							state = "TOURNAMENT_START_FINAL"
						} else {
							state = "TOURNAMENT_START_SEMI"
						}
					}
				} else if state == "TOURNAMENT_START_SEMI" {
					event := GameStart{
						TournamentEvent: TournamentEvent{
							Event: models.Event{
								Type: "TOURNAMENT_GAME",
							},
							Code: tournament.Id,
						},
						LobbyId: tournament.LobbiesSemi[0].Id,
					}
					evSemi1, _ := json.Marshal(&event)
					safeSend(tournament.LobbiesSemi[0].Sender.Send, evSemi1)
					safeSend(tournament.LobbiesSemi[0].Receiver.Send, evSemi1)
					event.LobbyId = tournament.LobbiesSemi[1].Id
					evSemi2, _ := json.Marshal(&event)
					safeSend(tournament.LobbiesSemi[1].Sender.Send, evSemi2)
					safeSend(tournament.LobbiesSemi[1].Receiver.Send, evSemi2)
					state = "TOURNAMENT_ON_SEMI"

					go func() {
						time.Sleep(100 * time.Millisecond)
						StartRoutine(h, tournament.LobbiesSemi[0])
						StartRoutine(h, tournament.LobbiesSemi[1])
					}()
				} else if state == "TOURNAMENT_ON_SEMI" {
					if tournament.Final[0] == 0 && tournament.LobbiesSemi[0].Game.State.Winner != 0 {
						fmt.Printf("semi1\n")
						tournament.Final[0] = tournament.LobbiesSemi[0].Game.State.Winner
						event.Semi1.Score[0] = uint8(tournament.LobbiesSemi[0].Game.State.Score.Player1)
						event.Semi1.Score[1] = uint8(tournament.LobbiesSemi[0].Game.State.Score.Player2)
						event.Semi1.IsFinished = true
						event.Final.Player1 = tournament.LobbiesSemi[0].Game.State.Winner
						jsonData, _ := json.Marshal(&event)
						SendDataToPlayers(tournament, jsonData)
					}
					if tournament.Final[1] == 0 && tournament.LobbiesSemi[1].Game.State.Winner != 0 {
						fmt.Printf("semi2\n")
						tournament.Final[1] = tournament.LobbiesSemi[1].Game.State.Winner
						event.Semi2.Score[0] = uint8(tournament.LobbiesSemi[1].Game.State.Score.Player1)
						event.Semi2.Score[1] = uint8(tournament.LobbiesSemi[1].Game.State.Score.Player2)
						event.Semi2.IsFinished = true
						event.Final.Player2 = tournament.LobbiesSemi[1].Game.State.Winner
						jsonData, _ := json.Marshal(&event)
						SendDataToPlayers(tournament, jsonData)
					}
					if tournament.Final[0] != 0 && tournament.Final[1] != 0 {
						fmt.Printf("TIMER_FINAL ON SEMI: %+v\n", tournament.Final)
						sec = 3
						close(tournament.LobbiesSemi[0].Destroy)
						close(tournament.LobbiesSemi[1].Destroy)
						state = "TIMER_FINAL"
					}
				} else if state == "TOURNAMENT_START_FINAL" {
					game := GameStart{
						TournamentEvent: TournamentEvent{
							Event: models.Event{
								Type: "TOURNAMENT_GAME",
							},
							Code: tournament.Id,
						},
						LobbyId: tournament.LobbyFinal.Id,
					}
					evFinal, _ := json.Marshal(&game)
					fmt.Printf("%+v\n", string(evFinal))
					safeSend(tournament.LobbyFinal.Sender.Send, evFinal)
					safeSend(tournament.LobbyFinal.Receiver.Send, evFinal)

					state = "TOURNAMENT_ON_FINAL"
					go func() {
						time.Sleep(100 * time.Millisecond)
						StartRoutine(h, tournament.LobbyFinal)
					}()
				} else if state == "TOURNAMENT_ON_FINAL" {
					if tournament.LobbyFinal.Game.State.Winner != 0 {
						fmt.Printf("ITS FINISH -> Userid %d won\n", tournament.LobbyFinal.Game.State.Winner)
						return
					}
				}
			}
		}
	}()

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
