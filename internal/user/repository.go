package user

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

func (r *Repository) Create(ctx context.Context, username, email, passwordHash string, createdAt time.Time) (int64, error) {
	query := `
			INSERT INTO users (username, email, password_hash, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id;
	`
	var id int64
	err := r.db.QueryRowContext(ctx, query, username, email, passwordHash, createdAt).Scan(&id)
	return id, err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
			SELECT id, username, email, password_hash, created_at 
			FROM users
			WHERE email = $1; 
	`

	user := &User{}
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
