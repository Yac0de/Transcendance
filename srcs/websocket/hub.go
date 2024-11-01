package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Hub struct {
	clients    map[uint64]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Message struct {
	SenderID   uint64    `json:"senderId"`
	ReceiverID uint64    `json:"receiverId"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[uint64]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func saveMessage(event Event) error {
	log.Printf("Message from %d to %d: %v", event.SenderID, event.ReceiverID, event.Data)
	message := Message{
		SenderID:   event.SenderID,
		ReceiverID: event.ReceiverID,
		Content:    event.Data,
		CreatedAt:  time.Now(),
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://backend:4000/conversation/add", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	return nil
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Id] = client
		case client := <-h.unregister:
			log.Printf("unregister client %d", client.Id)
			if _, ok := h.clients[client.Id]; ok {
				delete(h.clients, client.Id)
				close(client.Send)
			}
		case message := <-h.broadcast:
			var event Event
			if err := json.Unmarshal(message, &event); err != nil {
				log.Printf("error parsing message: %v", err)
			}
			for id := range h.clients {
				if id == event.SenderID || id == event.ReceiverID {
					select {
					case h.clients[id].Send <- message:
					default:
						close(h.clients[id].Send)
						delete(h.clients, id)
					}
				}
			}
			err := saveMessage(event)
			if err != nil {
				log.Printf("err: %v\n", err)
			}
		}
	}
}
