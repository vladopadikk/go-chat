package messages

import (
	"context"
	"errors"
	"fmt"

	"github.com/vladopadikk/go-chat/internal/chat"
)

var ErrForbidden = errors.New("user is not a member of the chat")
var ErrChatNotFound = errors.New("chat not found")

type Service struct {
	repo     *Repository
	chatRepo *chat.Repository
}

func NewService(repo *Repository, chatRepo *chat.Repository) *Service {
	return &Service{repo, chatRepo}
}

func (s *Service) SendMessage(ctx context.Context, senderID int64, input SendMessageInput) (Message, error) {
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return Message{}, err
	}
	defer tx.Rollback()

	exists, err := s.chatRepo.ChatExists(ctx, input.ChatID)
	if err != nil {
		return Message{}, err
	}
	if !exists {
		return Message{}, ErrChatNotFound
	}

	isMember, err := s.chatRepo.IsUserInChat(ctx, input.ChatID, senderID)
	if err != nil {
		return Message{}, err
	}
	if !isMember {
		return Message{}, ErrForbidden
	}

	msg, err := s.repo.Create(ctx, tx, Message{
		ChatID:   input.ChatID,
		SenderID: senderID,
		Content:  input.Content,
	})
	if err != nil {
		return Message{}, err
	}

	if err := tx.Commit(); err != nil {
		return Message{}, err
	}

	return msg, nil
}

func (s *Service) GetMessages(ctx context.Context, chatID, userID int64, limit, offset int) (MessageListResponse, error) {
	isMember, err := s.chatRepo.IsUserInChat(ctx, chatID, userID)
	if err != nil {
		return MessageListResponse{}, fmt.Errorf("db error: %w", err)
	}
	if !isMember {
		return MessageListResponse{}, ErrForbidden
	}

	msgs, err := s.repo.GetMsgByChatID(ctx, s.repo.db, chatID, limit, offset)
	if err != nil {
		return MessageListResponse{}, fmt.Errorf("db error: %w", err)
	}

	return MessageListResponse{
		Messages: msgs,
	}, nil
}
