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

func (s *Service) Register(ctx context.Context, input UserInput) (UserResponse, error) {
	_, err := s.repo.GetByEmail(ctx, input.Email)

	if err == nil {
		return UserResponse{}, ErrEmailExists
	}

	if err != sql.ErrNoRows {
		return UserResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, err
	}
	createdAt := time.Now()

	id, err := s.repo.Create(ctx, input.Username, input.Email, string(hashedPass), createdAt)
	if err != nil {
		return UserResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	return UserResponse{
		ID:       id,
		Username: input.Username,
		Email:    input.Email,
	}, nil
}
