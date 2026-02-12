package chat

import (
	"context"
	"gochat-server/internal/db"
	"gochat-server/internal/models"
	"time"
)

type Hub struct {
	rooms      map[string]map[*Client]bool
	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
	storage    *db.Storage
}

func NewHub(s *db.Storage) *Hub {
	return &Hub{
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
		storage:    s,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// when registering, add the client to a specific room.
			if h.rooms[client.chatID] == nil {
				h.rooms[client.chatID] = make(map[*Client]bool)
			}
			h.rooms[client.chatID][client] = true

		case client := <-h.unregister:
			if connections, ok := h.rooms[client.chatID]; ok {
				if _, ok := connections[client]; ok {
					delete(connections, client)
					close(client.send)
					if len(connections) == 0 {
						delete(h.rooms, client.chatID)
					}
				}
			}

		case msg := <-h.broadcast:

			// SHOULD BE REWRITTEN INTO A SEPARATE GOROUTINE and prob delete hub.storage - 15,18 lines
			go func(msg models.Message) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				h.storage.SaveMessage(ctx, msg)
			}(msg)
			// broadcast to chat members
			connections := h.rooms[msg.ChatID]
			for client := range connections {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(connections, client)
				}
			}
		}
	}
}
