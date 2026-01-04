package messages

import (
	"context"
	"database/sql"

	"github.com/vladopadikk/go-chat/internal/database"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(ctx context.Context, exec database.Executor, msg Message) (Message, error) {
	query := `
		INSERT INTO messages (chat_id, sender_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, chat_id, sender_id, content, created_at
	`
	var message Message

	err := exec.QueryRowContext(ctx, query, msg.ChatID, msg.SenderID, msg.Content).Scan(
		&message.ID,
		&message.ChatID,
		&message.SenderID,
		&message.Content,
		&message.CreatedAt,
	)
	return message, err
}

func (r *Repository) GetMsgByChatID(ctx context.Context, exec database.Executor, chatID int64, limit int, offset int) ([]MessageResponse, error) {
	query := `
		SELECT m.id, m.chat_id, m.sender_id, m.content, m.created_at 
		FROM messages m
		WHERE m.chat_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`
	rows, err := exec.QueryContext(ctx, query, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []MessageResponse

	for rows.Next() {
		var msg MessageResponse
		if err := rows.Scan(
			&msg.ID,
			&msg.ChatID,
			&msg.SenderID,
			&msg.Content,
			&msg.CreatedAt,
		); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, err

}
