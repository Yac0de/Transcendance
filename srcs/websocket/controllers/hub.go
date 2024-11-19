package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
	"websocket/models"
)

type Hub struct {
	Clients    map[uint64]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Lobbies    map[uuid.UUID]*Lobby
}

// type GameMessage struct {
// 	Type     string          `json:"type"`
// 	Command  string          `json:"command,omitempty"`
// 	PlayerID uint64          `json:"playerId,omitempty"`
// 	State    *GameState `json:"state,omitempty"`
// }

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

func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.Clients[client.Id]; !ok {
		return
	}

	target := h.Clients[client.Id]
	for id := range h.Lobbies {
		if h.Lobbies[id].Sender == target || h.Lobbies[id].Receiver == target {
			LobbyClientHasLeft(h, h.Lobbies[id].Id)
		}
	}
	go func() {
		time.Sleep(200 * time.Millisecond)
		NotifyClients(h, client.Id, "USER_DISCONNECTED")
	}()
	delete(h.Clients, client.Id)
	close(client.Send)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.Id] = client
			SendOnlineUsersToClient(h, client)
			NotifyClients(h, client.Id, "NEW_CONNECTION")
		case client := <-h.Unregister:
			h.RemoveClient(client)
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
			case event.Type == "GAME_EVENT":
				handleGameMessage(h, message)
			default:
				fmt.Printf("models.Event not handled: %+v\n", event)
			}
		}
	}
}

func safeSend(ch chan []byte, message []byte) {
	defer func() {
		if recover() != nil {
			fmt.Println("Attempted to send on a closed channel")
		}
	}()
	select {
	case ch <- message:
	default:
		fmt.Println("Channel is not ready to receive or closed")
	}
}

func safeClose(ch chan struct{}) {
	defer func() {
		if recover() != nil {
			fmt.Println("Attempted to close an already closed channel")
		}
	}()
	close(ch)
}
