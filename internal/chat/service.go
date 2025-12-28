package chat

import (
	"context"
	"database/sql"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) CreatePrivateChat(ctx context.Context, userID int64, createPrivateChatIn CreatePrivateChatInput) (Chat, error) {
	privateChat, err := s.repo.FindPrivateChatBetweenUsers(ctx, userID, createPrivateChatIn.UserID)
	if err != nil && err != sql.ErrNoRows {
		return Chat{}, fmt.Errorf("db error: %w", err)
	}

	if err == sql.ErrNoRows {
		privateChat, err = s.repo.CreateChat(ctx, createPrivateChatIn.Type, "")
		if err != nil {
			return Chat{}, fmt.Errorf("db error: %w", err)
		}
	}

	return privateChat, nil
}
