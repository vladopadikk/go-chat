package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/vladopadikk/go-chat/internal/config"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password")

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{repo, cfg}

}

func (s *Service) Login(ctx context.Context, loginIn LoginInput) (TokenResponse, error) {
	u, err := s.repo.GetByEmail(ctx, loginIn.Email)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("db error: %w", err)
	}

	if u == nil {
		return TokenResponse{}, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(loginIn.Password))
	if err != nil {
		return TokenResponse{}, ErrInvalidPassword
	}

	accessToken, err := GenerateAccessToken(u.ID, s.cfg.JWTSecret)
	if err != nil {
		return TokenResponse{}, err
	}
	refreshToken, err := GenerateRefreshToken(u.ID, s.cfg.JWTSecret)
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
