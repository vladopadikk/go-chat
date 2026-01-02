package messages

import "time"

type Message struct {
	ID        int64
	ChatID    int64
	SenderID  int64
	Content   string
	CreatedAt time.Time
}

type SendMessageInput struct {
	ChatID  int64  `json:"chat_id"`
	Content string `json:"content"`
}

type MessageResponse struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
