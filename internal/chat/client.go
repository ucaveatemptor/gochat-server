package chat

import (
	"gochat-server/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan models.Message
	chatID string
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg models.Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			break
		}
		msg.ChatID = c.chatID
		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			break
		}
	}
	c.conn.Close()
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chatId")
	if chatID == "" {
		http.Error(w, "chatId is required", http.StatusBadRequest)
		return
	}

	// upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// client init with certain chat
	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan models.Message, 256),
		chatID: chatID,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
