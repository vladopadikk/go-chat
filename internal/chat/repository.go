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

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *Repository) CreateChat(ctx context.Context, exec Executor, typ, name string) (Chat, error) {
	query := `
		INSERT INTO chats (type, name)
		VALUES ($1, $2)
		RETURNING id, type, created_at
	`
	var chat Chat
	err := exec.QueryRowContext(ctx, query, typ, name).Scan(&chat.ID, &chat.Type, &chat.CreatedAt)
	return chat, err
}

func (r *Repository) AddMember(ctx context.Context, exec Executor, chatID, userID int64, joinedAt time.Time) error {
	query := `
		INSERT INTO chat_members (chat_id, user_id, joined_at)
		VALUES ($1, $2, $3)
	`
	_, err := exec.ExecContext(ctx, query, chatID, userID, joinedAt)
	return err
}

func (r *Repository) GetChatsByUserID(ctx context.Context, userID int64) ([]ChatResponse, error) {
	query := `
		SELECT c.id, c.type, c.created_at 
		FROM chats c
		JOIN chat_members cm ON cm.chat_id = c.id
		WHERE cm.user_id = $1;
	`
	var chats []ChatResponse
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat ChatResponse
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
		JOIN chat_members cm2 ON cm2.chat_id = c.id
		WHERE c.type = 'private' AND cm1.user_id = $1 AND cm2.user_id = $2;
	`
	var chat Chat
	err := r.db.QueryRowContext(ctx, query, userA, userB).Scan(&chat.ID, &chat.Type, &chat.CreatedAt)
	return chat, err
}
