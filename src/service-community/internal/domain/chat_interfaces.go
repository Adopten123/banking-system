package domain

import (
	"context"
	"github.com/google/uuid"
)

type ChatRepository interface {
	CreateMessage(ctx context.Context, input CreateMessageInput) (*Message, error)
	GetChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int32) ([]Message, error)
}

type ChatService interface {
	SendMessage(ctx context.Context, input CreateMessageInput) (*Message, error)
	GetHistory(ctx context.Context, chatID, userID uuid.UUID, limit, offset int32) ([]Message, error)
}
