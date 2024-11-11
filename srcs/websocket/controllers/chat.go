package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"websocket/models"
)

func SaveMessageToDB(event models.MessageEvent) error {
	var message struct {
		SenderID   uint64    `json:"senderId"`
		ReceiverID uint64    `json:"receiverId"`
		Content    string    `json:"content"`
		CreatedAt  time.Time `json:"createdAt"`
	}

	message.SenderID = event.SenderID
	message.ReceiverID = event.ReceiverID
	message.Content = event.Data
	message.CreatedAt = time.Now()

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
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	return nil
}

func CreateOnlineUsersEvent(clients map[uint64]*Client, clientId uint64) models.OnlineUsersEvent {
	UsersOnline := models.OnlineUsersEvent{}
	UsersOnline.Type = "ONLINE_USERS"
	for id := range clients {
		if id != clientId {
			UsersOnline.Users = append(UsersOnline.Users, id)
		}
	}
	return UsersOnline
}

func CreateUserStatusEvent(id uint64, event string) models.UserStatusEvent {
	return models.UserStatusEvent{
		Event: models.Event{
			Type: event,
		},
		User: id,
	}
}

func SendOnlineUsersToClient(h *Hub, client *Client) {
	message, _ := json.Marshal(CreateOnlineUsersEvent(h.Clients, client.Id))
	select {
	case h.Clients[client.Id].Send <- message:
	default:
		close(h.Clients[client.Id].Send)
		delete(h.Clients, client.Id)
	}
}

func NotifyClients(h *Hub, clientId uint64, event string) {
	message, _ := json.Marshal(CreateUserStatusEvent(clientId, event))
	fmt.Printf("Event: %+v\n", string(message))
	for id := range h.Clients {
		if id != clientId {
			select {
			case h.Clients[id].Send <- message:
			default:
				close(h.Clients[id].Send)
				delete(h.Clients, id)
			}
		}
	}

}

func HandleChatMessage(h *Hub, message []byte) {
	var event models.MessageEvent
	if err := json.Unmarshal(message, &event); err != nil {
		fmt.Printf("error parsing message: %v", err)
		return
	}
	receiver, exists := h.Clients[event.ReceiverID]
	if exists {
		receiver.Send <- message
	}
	sender, exists := h.Clients[event.SenderID]
	if exists {
		sender.Send <- message
	}
	err := SaveMessageToDB(event)
	if err != nil {
		fmt.Printf("Error on saving data in db: %v\n", err)
	}
}
