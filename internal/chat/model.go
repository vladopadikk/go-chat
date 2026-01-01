package chat

import "time"

type Chat struct {
	ID        int64
	Type      string
	CreatedAt time.Time
}

type ChatMember struct {
	ChatID   int64
	UserID   int64
	JoinedAt time.Time
}

type CreatePrivateChatInput struct {
	UserID int64 `json:"user_id"`
}

type CreateGroupChatInput struct {
	Name         string  `json:"name"`
	Participants []int64 `json:"participants"`
}

type ChatResponse struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatListResponse struct {
	UserID int64
	Chats  []ChatResponse
}
