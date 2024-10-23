package main

import (
	"encoding/json"
	"log"
)

type Hub struct {
	clients    map[uint64]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[uint64]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
		case client := <-h.unregister:
			log.Printf("unregister client %d", client.id)
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		case message := <-h.broadcast:
			var event Event
			log.Printf("received a message: %s", string(message))
			if err := json.Unmarshal(message, &event); err != nil {
				log.Printf("error parsing message: %v", err)
			}
			for id := range h.clients {
				if id == event.SenderID || id == event.ReceiverID {
					select {
					case h.clients[id].send <- message:
					default:
						close(h.clients[id].send)
						delete(h.clients, id)
					}
				}
			}
		}
	}
}
