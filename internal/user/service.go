package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrEmailExists = errors.New("email is already registered")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Register(ctx context.Context, user UserInput) (int64, error) {
	u, err := s.repo.GetByEmail(ctx, user.Email)

	if u != nil {
		return 0, ErrEmailExists
	}

	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	createdAt := time.Now()

	return s.repo.Create(ctx, user.Username, user.Email, string(hashedPass), createdAt)
}
