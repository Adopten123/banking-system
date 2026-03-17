package service

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/google/uuid"
)

func (s *chatService) SendMessage(ctx context.Context, input domain.CreateMessageInput) (*domain.Message, error) {
	if (input.Content == nil || *input.Content == "") && !input.IsTransfer {
		return nil, fmt.Errorf("message content cannot be empty")
	}
	return s.repo.CreateMessage(ctx, input)
}

func (s *chatService) GetHistory(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int32) ([]domain.Message, error) {
	// TODO: add chack is userID in chat

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	return s.repo.GetChatMessages(ctx, chatID, limit, offset)
}

