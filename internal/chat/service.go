package chat

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) CreatePrivateChat(ctx context.Context, userID int64, createPrivateChatIn CreatePrivateChatInput) (Chat, error) {
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return Chat{}, fmt.Errorf("transaction error: %w", err)
	}
	defer tx.Rollback()

	chat, err := s.repo.FindPrivateChatBetweenUsers(ctx, userID, createPrivateChatIn.UserID)
	if err != nil && err != sql.ErrNoRows {
		return Chat{}, fmt.Errorf("db error: %w", err)
	}

	if err == sql.ErrNoRows {
		chat, err = s.repo.CreateChat(ctx, tx, "private", "")
		if err != nil {
			return Chat{}, fmt.Errorf("db error: %w", err)
		}
	}

	joinedAt := time.Now()

	err = s.repo.AddMember(ctx, tx, chat.ID, userID, joinedAt)
	if err != nil {
		return Chat{}, fmt.Errorf("db error: %w", err)
	}
	err = s.repo.AddMember(ctx, tx, chat.ID, createPrivateChatIn.UserID, joinedAt)
	if err != nil {
		return Chat{}, fmt.Errorf("db error: %w", err)
	}

	return chat, tx.Commit()
}

func (s *Service) CreateGroupChat(ctx context.Context, userID int64, createGroupChatIn CreateGroupChatInput) (Chat, error) {
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return Chat{}, fmt.Errorf("transaction error: %w", err)
	}
	defer tx.Rollback()

	chat, err := s.repo.CreateChat(ctx, tx, "group", createGroupChatIn.Name)
	if err != nil {
		return Chat{}, fmt.Errorf("db error: %w", err)
	}

	joinedAt := time.Now()

	err = s.repo.AddMember(ctx, tx, chat.ID, userID, joinedAt)
	if err != nil {
		return Chat{}, fmt.Errorf("db error: %w", err)
	}

	for _, participant := range createGroupChatIn.Participants {
		err = s.repo.AddMember(ctx, tx, chat.ID, participant, joinedAt)
		if err != nil {
			return Chat{}, fmt.Errorf("db error: %w", err)
		}
	}

	return chat, tx.Commit()
}

func (s *Service) GetChatsList(ctx context.Context, userID int64) ([]Chat, error) {
	chatList, err := s.repo.GetChatsByUserID(ctx, userID)
	if err != nil {
		return []Chat{}, fmt.Errorf("db error: %w", err)
	}

	return chatList, nil
}
