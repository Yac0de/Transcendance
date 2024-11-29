package controllers

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"websocket/models"

	"github.com/google/uuid"
)

const (
	PointPauseTime = 1 * time.Second
	GameTickRate   = 16 * time.Millisecond
)

type LobbyTimestamps struct {
	Pregame   time.Time `json:"pregame"`
	GameStart time.Time `json:"start"`
	GameEnd   time.Time `json:"end"`
}

type Lobby struct {
	Id           uuid.UUID       `json:"id"`
	Sender       *Client         `json:"sender"`
	Receiver     *Client         `json:"receiver"`
	Timestamps   LobbyTimestamps `json:"timestamps"`
	Status       string          `json:"status"`
	PlayersReady [2]bool         `json:"playersReady"`
	Mutex        sync.Mutex      `json:"-"`
	Destroy      chan struct{}   `json:"-"`
	Game         *Game           `json:"game"`
	IsTournamentGame bool            `json:"isTournamentGame"`
	IsActive	 bool			 `json:"isActive"`
	IsGameMode	 bool			 `json:"isGameMode"`
}

type LobbyUserState struct {
	Id      uint64 `json:"id"`
	IsReady bool   `json:"isReady"`
}

type LobbyEvent struct {
	models.Event
	LobbyId  uuid.UUID      `json:"lobbyId"`
	UserId   uint64         `json:"userId"`
	Sender   LobbyUserState `json:"sender"`
	Receiver LobbyUserState `json:"receiver"`
	IsTournamentGame bool           `json:"isTournamentGame"`
	IsGameMode	 bool			 `json:"isGameMode"`
}

type LobbyErrorEvent struct {
	models.Event
	LobbyId uuid.UUID `json:"lobbyId"`
	Error   string    `json:"error"`
}

func HandleLobby(h *Hub, event string, data []byte) {
	var request LobbyEvent
	if err := json.Unmarshal(data, &request); err != nil {
		fmt.Printf("Impossible to parse LobbyCreationRequest type: ", err.Error())
		return
	}
	switch event {
	case "LOBBY_INVITATION_TO_FRIEND":
		LobbyInvitation(h, request)
	case "LOBBY_ACCEPT_FROM_FRIEND":
		LobbyCreation(h, request)
	case "LOBBY_DENY_FROM_FRIEND":
		LobbyDenied(h, request)
	case "LOBBY_TERMINATE":
		LobbyTerminate(h, request)
	case "LOBBY_SPECIAL_MODE_TOGGLED":
		UpdateSpecialMode(h, request)
	case "LOBBY_PLAYER_READY_STATUS":
		LobbyUpdatePlayerStatus(h, request)
	case "LOBBY_PLAYER_UNREADY_STATUS":
		LobbyUpdatePlayerStatus(h, request)
	}
}

func (l *Lobby) ArePlayersReachable() bool {
	if l.Sender == nil || l.Receiver == nil {
		return false
	}
	return true
}

func LobbyClientHasLeft(h *Hub, lobbyId uuid.UUID) {
	lobby := h.Lobbies[lobbyId]
	error := LobbyErrorEvent{
		Event: models.Event{
			Type: "LOBBY_DESTROYED",
		},
		LobbyId: lobby.Id,
		Error:   "A player has left the lobby",
	}

	errorJson, err := json.Marshal(&error)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyErrorEvent type: ", err.Error())
		return
	}

	safeSend(lobby.Sender.Send, errorJson)
	safeSend(lobby.Receiver.Send, errorJson)

	if lobby.Destroy != nil {
		fmt.Printf("lobby destroyed lobby %+v\n", lobby.Id)
		safeClose(lobby.Destroy)
	}
	delete(h.Lobbies, lobby.Id)
}

func LobbyInvitation(h *Hub, request LobbyEvent) {
	lobbyId := uuid.New()
	request.LobbyId = lobbyId

	senderJson, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}
	safeSend(h.Clients[request.Sender.Id].Send, senderJson)

	request.Type = "LOBBY_INVITATION_FROM_FRIEND"
	receiverJson, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}
	safeSend(h.Clients[request.Receiver.Id].Send, receiverJson)
}

func LobbyCreation(h *Hub, request LobbyEvent) {
	lobby, err := NewLobby(h, request)
	if err != nil {
		fmt.Printf("Lobby creation failed : ", err.Error())
		return
	}
	h.Lobbies[lobby.Id] = lobby

	response := LobbyEvent{
		Event: models.Event{
			Type: "LOBBY_CREATED",
		},
		Sender: LobbyUserState{
			Id:      lobby.Sender.Id,
			IsReady: false,
		},
		Receiver: LobbyUserState{
			Id:      lobby.Receiver.Id,
			IsReady: false,
		},
		LobbyId:          lobby.Id,
		IsTournamentGame: false,
	}
	jsonData, err := json.Marshal(&response)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}

	safeSend(lobby.Sender.Send, jsonData)
	safeSend(lobby.Receiver.Send, jsonData)
}

func LobbyDenied(h *Hub, request LobbyEvent) {
	request.Type = "LOBBY_DENIED"
	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}

	safeSend(h.Clients[request.Sender.Id].Send, jsonData)
	safeSend(h.Clients[request.Receiver.Id].Send, jsonData)
}

func LobbyTerminate(h *Hub, request LobbyEvent) {
	lobby, exists := h.Lobbies[request.LobbyId]
	if !exists {
		fmt.Printf("Trying to remove a lobby doesn't exists: %+v\n", request)
		return
	}
	request.Type = "LOBBY_DESTROYED"
	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}

	safeSend(lobby.Sender.Send, jsonData)
	safeSend(lobby.Receiver.Send, jsonData)

	if lobby.Destroy != nil {
		safeClose(lobby.Destroy)
	}
	delete(h.Lobbies, request.LobbyId)
}

func LobbyUpdatePlayerStatus(h *Hub, request LobbyEvent) {
	lobby, exists := h.Lobbies[request.LobbyId]
	if !exists {
		fmt.Printf("Lobby %s not found\n", request.LobbyId)
		return
	}

	if lobby.ArePlayersReachable() == false {
		LobbyClientHasLeft(h, lobby.Id)
		return
	}

	isReady := false
	if request.Type == "LOBBY_PLAYER_READY_STATUS" {
		isReady = true
	}

	fmt.Printf("Request: \n", request)
	if request.UserId == lobby.Sender.Id {
		lobby.PlayersReady[0] = isReady
	} else if request.UserId == lobby.Receiver.Id {
		lobby.PlayersReady[1] = isReady
	}

	request.Sender = LobbyUserState{
		Id:      lobby.Sender.Id,
		IsReady: lobby.PlayersReady[0],
	}
	request.Receiver = LobbyUserState{
		Id:      lobby.Receiver.Id,
		IsReady: lobby.PlayersReady[1],
	}

	request.Type = "LOBBY_PLAYER_STATUS"
	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}

	safeSend(lobby.Sender.Send, jsonData)
	safeSend(lobby.Receiver.Send, jsonData)
	if lobby.PlayersReady[0] && lobby.PlayersReady[1] {
		go func() {
			time.Sleep(10 * time.Millisecond)
			//fmt.Printf("Lobby game mode %v\n\n\n", lobby.IsGameMode)
			StartRoutine(h, lobby)
		}()
		return
	}
}

func UpdateSpecialMode(h *Hub, request LobbyEvent) {
    // Récupérer le lobby à partir de son ID
    lobby, exists := h.Lobbies[request.LobbyId]
    if !exists {
        fmt.Printf("Lobby not found: %s\n", request.LobbyId)
        return
    }

    // Mettre à jour le mode spécial pour ce lobby
	fmt.Printf("GAME MODE RECEIVED\n", request.IsGameMode)
    lobby.IsGameMode = request.IsGameMode

    // Créer un événement et le transmettre aux clients du lobby
    event := LobbyEvent{
        Event: models.Event{
            Type: "LOBBY_SPECIAL_MODE_TOGGLED",
        },
        LobbyId:   lobby.Id,
        IsGameMode: request.IsGameMode,
    }

    // Serialiser l'événement et l'envoyer aux clients
    senderJson, err := json.Marshal(&event)
    if err != nil {
        fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
        return
    }

    safeSend(lobby.Sender.Send, senderJson)
    safeSend(lobby.Receiver.Send, senderJson)
}

func StartRoutine(h *Hub, lobby *Lobby) {
	lobby.Timestamps.Pregame = time.Now()
	lobby.Destroy = make(chan struct{})
	lobby.Game = NewGame(lobby.Sender.Id, lobby.Receiver.Id)
	gameTicker := time.NewTicker(GameTickRate)

	gameStart := models.Event{
		Type: "GAME_START",
	}

	dataJson, err := json.Marshal(gameStart)
	if err != nil {
		fmt.Printf("Impossible to parse GameStart type: ", err.Error())
		return
	}

	safeSend(lobby.Sender.Send, dataJson)
	safeSend(lobby.Receiver.Send, dataJson)

	lobby.Game.resetBall()
	go func() {
		for {
			select {
			case <-lobby.Destroy:
				gameTicker.Stop()
				return
			case <-gameTicker.C:
				if lobby.Game != nil && lobby.Game.State.IsActive {
					lobby.Game.Update()
					evt := GameEvent{
						Event: models.Event{
							Type: "GAME_EVENT",
						},
						LobbyId:          lobby.Id,
						State:            lobby.Game.State,
						Player1Id:        lobby.Sender.Id,
						Player2Id:        lobby.Receiver.Id,
						IsTournamentGame: lobby.IsTournamentGame,
					}
					stateJson, _ := json.Marshal(evt)
					safeSend(lobby.Sender.Send, stateJson)
					safeSend(lobby.Receiver.Send, stateJson)
				} else if lobby.Game.State.IsActive == false && lobby.Game.State.Winner != 0 {
					evt := GameEvent{
						Event: models.Event{
							Type: "GAME_FINISHED",
						},

						LobbyId:          lobby.Id,
						State:            lobby.Game.State,
						Player1Id:        lobby.Sender.Id,
						Player2Id:        lobby.Receiver.Id,
						IsTournamentGame: lobby.IsTournamentGame,
					}
					stateJson, _ := json.Marshal(evt)
					safeSend(lobby.Sender.Send, stateJson)
					safeSend(lobby.Receiver.Send, stateJson)
					return
				}
			}
		}
	}()
}

func (lobby *Lobby) DispatchTimer(timeLeft time.Duration) {
	if lobby == nil {
		return
	}
	var RemainingTime struct {
		models.Event
		Time int `json:"remainingSecondsToStart"`
	}
	RemainingTime.Type = "LOBBY_PREGAME_REMAINING_TIME"
	RemainingTime.Time = int(timeLeft.Seconds())

	jsonData, err := json.Marshal(&RemainingTime)
	if err != nil {
		fmt.Printf("Impossible to parse RemainingTime type: ", err.Error())
		return
	}
	safeSend(lobby.Sender.Send, jsonData)
	safeSend(lobby.Receiver.Send, jsonData)
}

func NewLobby(h *Hub, request LobbyEvent) (*Lobby, error) {
	sender := h.Clients[request.Sender.Id]
	if sender == nil {
		return nil, fmt.Errorf("Sender.Id doesn't exists %d", request.Sender.Id)
	}

	receiver := h.Clients[request.Receiver.Id]
	if receiver == nil {
		return nil, fmt.Errorf("Receiver.Id doesn't exists %d", request.Receiver.Id)
	}
	newSession := &Lobby{
		Id:           request.LobbyId,
		Sender:       sender,
		Receiver:     receiver,
		Timestamps:   LobbyTimestamps{},
		Status:       "LOBBY_CREATION",
		PlayersReady: [2]bool{false, false},
	}
	return newSession, nil
}

func indexOf(element uint64, data []uint64) int {
	for i, v := range data {
		if v == element {
			return i
		}
	}
	return -1
}
