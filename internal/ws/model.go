package ws

import (
	"encoding/json"
	"time"
)

const (
	WSMessageTypeSendMessage = "send_message"
	WSMessageTypeNewMessage  = "new_message"
	WSMessageTypeError       = "error"
)

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessagePayload struct {
	ChatID  int64  `json:"chat_id"`
	Content string `json:"content"`
}

type NewMessagePayload struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}
