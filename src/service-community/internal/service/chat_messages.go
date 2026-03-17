package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/google/uuid"
)

func (s *chatService) SendMessage(ctx context.Context, input domain.CreateMessageInput) (*domain.Message, error) {
	if (input.Content == nil || *input.Content == "") && !input.IsTransfer {
		return nil, fmt.Errorf("message content cannot be empty")
	}

	msg, err := s.repo.CreateMessage(ctx, input)
	if err != nil {
		return nil, err
	}

	if msg.IsTransfer && msg.TransferAmount != nil && msg.IdempotencyKey != nil {
		event := domain.TransferCreatedEvent{
			MessageID:      msg.ID,
			ChatID:         msg.ChatID,
			SenderID:       msg.SenderID,
			Amount:         *msg.TransferAmount,
			Currency:       *msg.TransferCurrency,
			IdempotencyKey: *msg.IdempotencyKey,
		}

		go func() {
			err := s.publisher.PublishTransferCreated(context.Background(), event)
			if err != nil {
				log.Printf("CRITICAL: Failed to publish transfer event for msg %d: %v", msg.ID, err)
			}
		}()
	}

	return msg, nil
}

func (s *chatService) GetHistory(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int32) ([]domain.Message, error) {
	// TODO: add chack is userID in chat

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	return s.repo.GetChatMessages(ctx, chatID, limit, offset)
}

func (s *chatService) GetChatMemberIDs(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	return s.repo.GetChatMemberIDs(ctx, chatID)
}