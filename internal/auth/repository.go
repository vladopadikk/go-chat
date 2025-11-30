package auth

import (
	"context"
	"database/sql"

	"github.com/vladopadikk/go-chat/internal/user"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
			SELECT id, username, email, password_hash, created_at 
			FROM users
			WHERE email = $1; 
	`

	user := &user.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}
