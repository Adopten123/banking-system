package domain

import (
	"context"
	"github.com/google/uuid"
)

type ChatRepository interface {
	CreateMessage(ctx context.Context, input CreateMessageInput) (*Message, error)
	GetChatMessages(ctx context.Context, chatID int64, limit, offset int) ([]Message, error)
}

type ChatService interface {
	SendMessage(ctx context.Context, input CreateMessageInput) (*Message, error)
	GetHistory(ctx context.Context, chatID int64, userID uuid.UUID, limit, offset int) ([]Message, error)
}
