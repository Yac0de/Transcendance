package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"websocket/models"

	"github.com/google/uuid"
)

type Hub struct {
	Clients    map[uint64]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Lobbies    map[uuid.UUID]*Lobby
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint64]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Lobbies:    make(map[uuid.UUID]*Lobby),
	}
}

func (h *Hub) GetEventType(message []byte) (models.Event, error) {
	var evt models.Event
	if err := json.Unmarshal(message, &evt); err != nil {
		return models.Event{}, err
	}
	return evt, nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.Id] = client
			SendOnlineUsersToClient(h, client)
			NotifyClients(h, client.Id, "NEW_CONNECTION")
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.Id]; ok {
				delete(h.Clients, client.Id)
				close(client.Send)
				NotifyClients(h, client.Id, "USER_DISCONNECTED")
			}
		case message := <-h.Broadcast:
			event, err := h.GetEventType(message)
			if err != nil {
				fmt.Printf("Hub.broadcast error on event cast: %s | error: %v\n", string(message), err)
			}
			switch {
			case event.Type == "CHAT":
				HandleChatMessage(h, message)
			case strings.HasPrefix(event.Type, "LOBBY_"):
				HandleLobby(h, event.Type, message)
			default:
				fmt.Printf("models.Event not handled: %+v\n", event)
			}
		}
	}
}
