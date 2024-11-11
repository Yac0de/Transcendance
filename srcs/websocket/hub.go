package main

import (
	"encoding/json"
	"fmt"
)

type Hub struct {
	clients    map[uint64]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint64]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) GetEventType(message []byte) (Event, error) {
	var evt Event
	if err := json.Unmarshal(message, &evt); err != nil {
		return Event{}, err
	}
	return evt, nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Id] = client
			SendOnlineUsersToClient(h, client)
			NotifyClients(h, client.Id, "NEW_CONNECTION")
		case client := <-h.unregister:
			if _, ok := h.clients[client.Id]; ok {
				delete(h.clients, client.Id)
				close(client.Send)
				NotifyClients(h, client.Id, "USER_DISCONNECTED")
			}
		case message := <-h.broadcast:
			event, err := h.GetEventType(message)
			if err != nil {
				fmt.Printf("Hub.broadcast error on event cast: %s | error: %v\n", string(message), err)
			}
			switch event.Type {
			case "CHAT":
				HandleChatMessage(h, message)
			default:
				fmt.Printf("Event not handled: %+v\n", event)
			}
		}
	}
}
