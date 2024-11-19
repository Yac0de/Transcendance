package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
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

// func (h *Hub) handleGameMessage(message []byte) {
// 	var gameMsg GameMessage
// 	if err := json.Unmarshal(message, &gameMsg); err != nil {
// 		log.Printf("Error parsing game message: %v", err)
// 		return
// 	}

// 	// Ajout de ce log pour le debug
// 	log.Printf("Received game message: %+v", gameMsg) // AJOUTÉ

// 	switch gameMsg.Type {
// 	case "GAME_START":
// 		availablePlayers := make([]uint64, 0)
// 		for clientID := range h.Clients {
// 			availablePlayers = append(availablePlayers, clientID)
// 			if len(availablePlayers) == 2 {
// 				break
// 			}
// 		}

// 		if len(availablePlayers) < 2 {
// 			log.Println("Not enough players to start a game")
// 			return
// 		}

// 		newGame := NewGame(availablePlayers[0], availablePlayers[1])
// 		newGame.State.IsActive = true // AJOUTÉ: Activons explicitement le jeu
// 		h.games[newGame.ID] = newGame

// 		// Notify players that game has started
// 		startMsg := GameMessage{
// 			Type:  "GAME_START",
// 			State: newGame.State,
// 		}
// 		msgBytes, _ := json.Marshal(startMsg)
// 		for _, playerID := range availablePlayers {
// 			if client, ok := h.clients[playerID]; ok {
// 				// AJOUTÉ: Log pour le debug
// 				log.Printf("Sending start message to player: %d", playerID)
// 				client.Send <- msgBytes
// 			}
// 		}

// 		go h.runGame(newGame)

// 	case "GAME_COMMAND":
// 		for _, g := range h.Games {
// 			if g.Player1.ID == gameMsg.PlayerID || g.Player2.ID == gameMsg.PlayerID {
// 				g.HandleCommand(GameCommand{
// 					PlayerID: gameMsg.PlayerID,
// 					Command:  gameMsg.Command,
// 				})
// 				// AJOUTÉ: Log pour le debug
// 				log.Printf("Handled command %s for player %d", gameMsg.Command, gameMsg.PlayerID)
// 			}
// 		}
// 	}
// }

// func (h *Hub) runGame(g *game.Game) {
// 	ticker := time.NewTicker(16 * time.Millisecond)
// 	defer ticker.Stop()

// 	// AJOUTÉ: Log pour le debug
// 	log.Printf("Starting game loop. Initial state: active=%v", g.State.IsActive)

// 	for range ticker.C {
// 		if !g.State.IsActive {
// 			// AJOUTÉ: Log pour le debug
// 			log.Println("Game became inactive, stopping game loop")
// 			return
// 		}

// 		g.Update()

// 		stateMsg := GameMessage{
// 			Type:  "GAME_UPDATE",
// 			State: g.State,
// 		}
// 		msgBytes, _ := json.Marshal(stateMsg)

// 		if client1, ok := h.clients[g.Player1.ID]; ok {
// 			client1.Send <- msgBytes
// 		}
// 		if client2, ok := h.clients[g.Player2.ID]; ok {
// 			client2.Send <- msgBytes
// 		}

// 		if g.State.Winner != 0 {
// 			// AJOUTÉ: Log pour le debug
// 			log.Printf("Game ended. Winner: %d", g.State.Winner)
// 			delete(h.games, g.ID)
// 			return
// 		}
// 	}
// }
LobbyId  
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
