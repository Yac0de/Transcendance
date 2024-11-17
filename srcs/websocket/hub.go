package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"ws-back/game"
)

type GameMessage struct {
	Type     string          `json:"type"`
	Command  string          `json:"command,omitempty"`
	PlayerID uint64          `json:"playerId,omitempty"`
	State    *game.GameState `json:"state,omitempty"`
}

type Hub struct {
	clients    map[uint64]*Client
	games      map[string]*game.Game
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
		games:      make(map[string]*game.Game),
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

func (h *Hub) handleGameMessage(message []byte) {
	var gameMsg GameMessage
	if err := json.Unmarshal(message, &gameMsg); err != nil {
		log.Printf("Error parsing game message: %v", err)
		return
	}

	// Ajout de ce log pour le debug
	log.Printf("Received game message: %+v", gameMsg) // AJOUTÉ

	switch gameMsg.Type {
	case "GAME_START":
		availablePlayers := make([]uint64, 0)
		for clientID := range h.clients {
			availablePlayers = append(availablePlayers, clientID)
			if len(availablePlayers) == 2 {
				break
			}
		}

		if len(availablePlayers) < 2 {
			log.Println("Not enough players to start a game")
			return
		}

		newGame := game.NewGame(availablePlayers[0], availablePlayers[1])
		newGame.State.IsActive = true // AJOUTÉ: Activons explicitement le jeu
		h.games[newGame.ID] = newGame

		// Notify players that game has started
		startMsg := GameMessage{
			Type:  "GAME_START",
			State: newGame.State,
		}
		msgBytes, _ := json.Marshal(startMsg)
		for _, playerID := range availablePlayers {
			if client, ok := h.clients[playerID]; ok {
				// AJOUTÉ: Log pour le debug
				log.Printf("Sending start message to player: %d", playerID)
				client.Send <- msgBytes
			}
		}

		go h.runGame(newGame)

	case "GAME_COMMAND":
		for _, g := range h.games {
			if g.Player1.ID == gameMsg.PlayerID || g.Player2.ID == gameMsg.PlayerID {
				g.HandleCommand(game.GameCommand{
					PlayerID: gameMsg.PlayerID,
					Command:  gameMsg.Command,
				})
				// AJOUTÉ: Log pour le debug
				log.Printf("Handled command %s for player %d", gameMsg.Command, gameMsg.PlayerID)
			}
		}
	}
}

func (h *Hub) runGame(g *game.Game) {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	// AJOUTÉ: Log pour le debug
	log.Printf("Starting game loop. Initial state: active=%v", g.State.IsActive)

	for range ticker.C {
		if !g.State.IsActive {
			// AJOUTÉ: Log pour le debug
			log.Println("Game became inactive, stopping game loop")
			return
		}

		g.Update()

		stateMsg := GameMessage{
			Type:  "GAME_UPDATE",
			State: g.State,
		}
		msgBytes, _ := json.Marshal(stateMsg)

		if client1, ok := h.clients[g.Player1.ID]; ok {
			client1.Send <- msgBytes
		}
		if client2, ok := h.clients[g.Player2.ID]; ok {
			client2.Send <- msgBytes
		}

		if g.State.Winner != 0 {
			// AJOUTÉ: Log pour le debug
			log.Printf("Game ended. Winner: %d", g.State.Winner)
			delete(h.games, g.ID)
			return
		}
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
			// Try to parse as game message first
			DispatchMessage(h, message)
			var gameMsg GameMessage
			if err := json.Unmarshal(message, &gameMsg); err == nil && gameMsg.Type != "" {
				h.handleGameMessage(message)
				continue
			}
		}
	}
}
