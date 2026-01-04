package ws

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vladopadikk/go-chat/internal/messages"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait * 9 / 10
	maxMessageSize = 512 * 1024
)

type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan []byte
	userID         int64
	chats          []int64
	messageService *messages.Service
}

func NewClient(hub *Hub, conn *websocket.Conn, userID int64, chats []int64, messageService *messages.Service) *Client {
	return &Client{
		hub:            hub,
		conn:           conn,
		userID:         userID,
		chats:          chats,
		send:           make(chan []byte, 256),
		messageService: messageService,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		c.handleMessage(msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, msg)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(data []byte) {
	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.sendError("invalid message format")
		return
	}

	switch msg.Type {
	case WSMessageTypeSendMessage:
		c.handleSendMessage(msg.Payload)
	default:
		c.sendError("unknown message type")
	}
}

func (c *Client) handleSendMessage(payload json.RawMessage) {
	var input SendMessagePayload
	if err := json.Unmarshal(payload, &input); err != nil {
		c.sendError("invalid payload")
		return
	}

	msg, err := c.messageService.SendMessage(context.Background(), c.userID, messages.SendMessageInput(input))
	if err != nil {
		c.sendError(err.Error())
		return
	}

	out := NewMessagePayload{
		ID:        msg.ID,
		ChatID:    msg.ChatID,
		SenderID:  msg.SenderID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
	}

	payloadBytes, _ := json.Marshal(out)
	wsMsg, _ := json.Marshal(WSMessage{
		Type:    WSMessageTypeNewMessage,
		Payload: payloadBytes,
	})

	c.hub.broadcast <- Broadcast{
		ChatID: msg.ChatID,
		Data:   wsMsg,
	}
}

func (c *Client) sendError(text string) {
	payload, _ := json.Marshal(ErrorPayload{Message: text})
	msg, _ := json.Marshal(WSMessage{
		Type:    WSMessageTypeError,
		Payload: payload,
	})
	c.send <- msg
}
