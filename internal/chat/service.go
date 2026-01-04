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

func (s *Service) CreatePrivateChat(ctx context.Context, userID int64, createPrivateChatIn CreatePrivateChatInput) (ChatResponse, error) {
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("transaction error: %w", err)
	}
	defer tx.Rollback()

	chat, err := s.repo.FindPrivateChatBetweenUsers(ctx, s.repo.db, userID, createPrivateChatIn.UserID)
	if err != nil && err != sql.ErrNoRows {
		return ChatResponse{}, fmt.Errorf("db error: %w", err)
	}

	if err == sql.ErrNoRows {
		chat, err = s.repo.CreateChat(ctx, tx, "private", "")
		if err != nil {
			return ChatResponse{}, fmt.Errorf("db error: %w", err)
		}
		joinedAt := time.Now()

		err = s.repo.AddMember(ctx, tx, chat.ID, userID, joinedAt)
		if err != nil {
			return ChatResponse{}, fmt.Errorf("db error: %w", err)
		}
		err = s.repo.AddMember(ctx, tx, chat.ID, createPrivateChatIn.UserID, joinedAt)
		if err != nil {
			return ChatResponse{}, fmt.Errorf("db error: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return ChatResponse{}, fmt.Errorf("commit tx: %w", err)
	}

	return ChatResponse{
		ID:        chat.ID,
		Type:      chat.Type,
		CreatedAt: chat.CreatedAt,
	}, nil
}

func (s *Service) CreateGroupChat(ctx context.Context, userID int64, createGroupChatIn CreateGroupChatInput) (ChatResponse, error) {
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("transaction error: %w", err)
	}
	defer tx.Rollback()

	chat, err := s.repo.CreateChat(ctx, tx, "group", createGroupChatIn.Name)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("db error: %w", err)
	}

	joinedAt := time.Now()

	for _, participant := range createGroupChatIn.Participants {
		err = s.repo.AddMember(ctx, tx, chat.ID, participant, joinedAt)
		if err != nil {
			return ChatResponse{}, fmt.Errorf("db error: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return ChatResponse{}, fmt.Errorf("commit tx: %w", err)
	}

	return ChatResponse{
		ID:        chat.ID,
		Type:      chat.Type,
		CreatedAt: chat.CreatedAt,
	}, nil
}

func (s *Service) GetChatsList(ctx context.Context, userID int64) (ChatListResponse, error) {
	chatList, err := s.repo.GetChatsByUserID(ctx, s.repo.db, userID)
	if err != nil {
		return ChatListResponse{}, fmt.Errorf("db error: %w", err)
	}

	return ChatListResponse{
		Chats: chatList,
	}, nil
}
