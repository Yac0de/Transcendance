package controllers

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"websocket/models"

	"github.com/google/uuid"
)

type LobbyTimestamps struct {
	Request  time.Time `json:"request"`
	Response time.Time `json:"response"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

type Lobby struct {
	Id           uuid.UUID `json:"id"`
	Sender       *Client   `json:"sender"`
	Receiver     *Client   `json:"receiver"`
	PlayerJoined [2]uint64 `json:"playerJoined"`
	// Instance    *Game          `json:"instance"`
	Timestamps LobbyTimestamps `json:"timestamps"`
	Status     string          `json:"status"`
	Mutex      sync.Mutex      `json:"-"`
	Destroy    chan struct{}   `json:"-"`
}

type LobbyEvent struct {
	models.Event
	UserId     uint64    `json:"userId"`
	SenderID   uint64    `json:"senderId"`
	ReceiverID uint64    `json:"receiverId"`
	LobbyId    uuid.UUID `json:"lobbyId"`
	Status     string    `json:"status"`
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
	}
}

func LobbyInvitation(h *Hub, request LobbyEvent) {
	lobbyId := uuid.New()
	request.LobbyId = lobbyId

	senderJson, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}
	h.Clients[request.SenderID].Send <- senderJson

	request.Type = "LOBBY_INVITATION_FROM_FRIEND"
	receiverJson, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}
	h.Clients[request.ReceiverID].Send <- receiverJson
}

func LobbyTerminate(h *Hub, request LobbyEvent) {
	_, exists := h.Lobbies[request.LobbyId]
	if exists {
		// close Destroy field
		delete(h.Lobbies, request.LobbyId)
	}
	request.Type = "LOBBY_DESTROYED"
	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}
	if _, exists := h.Clients[request.SenderID]; exists {
		h.Clients[request.SenderID].Send <- jsonData
	}

	if _, exists := h.Clients[request.ReceiverID]; exists {
		h.Clients[request.ReceiverID].Send <- jsonData
	}
}

func LobbyDenied(h *Hub, request LobbyEvent) {
	request.Type = "LOBBY_DENIED"
	jsonData, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}
	h.Clients[request.SenderID].Send <- jsonData
	h.Clients[request.ReceiverID].Send <- jsonData
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
		LobbyId:    lobby.Id,
		SenderID:   lobby.Sender.Id,
		ReceiverID: lobby.Receiver.Id,
		Status:     "INITILIZED",
	}
	jsonData, err := json.Marshal(&response)
	if err != nil {
		fmt.Printf("Impossible to parse LobbyEvent type: ", err.Error())
		return
	}

	lobby.Sender.Send <- jsonData
	lobby.Receiver.Send <- jsonData
}

func NewLobby(h *Hub, request LobbyEvent) (*Lobby, error) {
	sender := h.Clients[request.SenderID]
	if sender == nil {
		return nil, fmt.Errorf("SenderID doesn't exists %d", request.SenderID)
	}

	receiver := h.Clients[request.ReceiverID]
	if receiver == nil {
		return nil, fmt.Errorf("ReceiverID doesn't exists %d", request.ReceiverID)
	}
	newSession := &Lobby{
		Id:       request.LobbyId,
		Sender:   sender,
		Receiver: receiver,
		// Instance: nil,
		Timestamps: LobbyTimestamps{
			Request: time.Now(),
		},
		Status: "LOBBY_CREATION",
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
