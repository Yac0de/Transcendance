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

func saveMessage(event MessageEvent) error {
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

func CreateOnlineUsersEvent(clients map[uint64]*Client, clientId uint64) OnlineUsersEvent {
	UsersOnline := OnlineUsersEvent{}
	UsersOnline.Type = "ONLINE_USERS"
	for id := range clients {
		if id != clientId {
			UsersOnline.Users = append(UsersOnline.Users, id)
		}
	}
	return UsersOnline
}

func CreateUserStatusEvent(id uint64, event string) UserStatusEvent {
	return UserStatusEvent{
		Event: Event{
			Type: event,
		},
		User: id,
	}
}

func SendOnlineUsersToClient(h *Hub, client *Client) {
	message, _ := json.Marshal(CreateOnlineUsersEvent(h.clients, client.Id))
	select {
	case h.clients[client.Id].Send <- message:
	default:
		close(h.clients[client.Id].Send)
		delete(h.clients, client.Id)
	}
}

func NotifyClients(h *Hub, clientId uint64, event string) {
	message, _ := json.Marshal(CreateUserStatusEvent(clientId, event))
	for id := range h.clients {
		if id != clientId {
			select {
			case h.clients[id].Send <- message:
			default:
				close(h.clients[id].Send)
				delete(h.clients, id)
			}
		}
	}

}

func DispatchMessage(h *Hub, message []byte) {
	var event MessageEvent
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

func (h *Hub) run() {
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
			DispatchMessage(h, message)
		}
	}
}
