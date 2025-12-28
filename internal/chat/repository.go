package chat

import (
	"context"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateChat(ctx context.Context, typ, name string) (Chat, error) {
	query := `
		INSERT INTO chats (type, name)
		VALUES ($1, $2, $3)
		RETURNING id, type, created_at
	`
	var chat Chat
	err := r.db.QueryRowContext(ctx, query, typ, name).Scan(&chat.ID, &chat.Type, &chat.CreatedAt)
	return chat, err
}

func (r *Repository) AddMember(ctx context.Context, chatID, userID int64, joinedAt time.Time) error {
	query := `
		INSERT INTO chat_members (chat_id, user_id, joined_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, chatID, userID, joinedAt)
	return err
}

func (r *Repository) GetChatsByUserID(ctx context.Context, userID int64) ([]Chat, error) {
	query := `
		SELECT c.id, c.type, c.created_at 
		FROM chats c
		JOIN chat_members cm ON cm.chat_id = c.id
		WHERE cm.user_id = $1;
	`
	var chats []Chat
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.Type, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, err
}

func (r *Repository) FindPrivateChatBetweenUsers(ctx context.Context, userA, userB int64) (Chat, error) {
	query := `
		SELECT c.id, c.type, c.created_at
		FROM chats c
		JOIN chat_members cm1 ON cm1.chat_id = c.id
		JOIN chat_members cm3 ON cm2.chat_id = c.id
		WHERE c.type == "private" AND cm.user_id = $1 AND cm.user_id = $2;
	`
	var chat Chat
	err := r.db.QueryRowContext(ctx, query, userA, userB).Scan(&chat.ID, &chat.Type, &chat.CreatedAt)
	return chat, err
}
