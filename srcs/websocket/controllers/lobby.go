package controllers

import (
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
	Id       uuid.UUID `json:"id"`
	Sender   *Client   `json:"sender"`
	Receiver *Client   `json:"receiver"`
	// Instance    *Game          `json:"instance"`
	Timestamps LobbyTimestamps `json:"timestamps"`
	Status     string          `json:"status"`
	Mutex      sync.Mutex      `json:"-"`
	Destroy    chan struct{}   `json:"-"`
}

type LobbyCreationRequest struct {
	models.Event
	SenderID   uint64 `json:"senderId"`
	ReceiverID uint64 `json:"receiverId"`
}

func HandleLobby(h *Hub, event string, data []byte) {

}

func NewLobby(h *Hub, Request LobbyCreationRequest) (*Lobby, error) {
	sender := h.Clients[Request.SenderID]
	if sender == nil {
		return nil, fmt.Errorf("SenderID doesn't exists %d", Request.SenderID)
	}

	receiver := h.Clients[Request.ReceiverID]
	if receiver == nil {
		return nil, fmt.Errorf("ReceiverID doesn't exists %d", Request.ReceiverID)
	}
	newSession := &Lobby{
		Id:       uuid.New(),
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
