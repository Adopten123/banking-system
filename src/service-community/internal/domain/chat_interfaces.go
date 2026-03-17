package domain

import (
	"context"
	"github.com/google/uuid"
)

type ChatRepository interface {
	CreateMessage(ctx context.Context, input CreateMessageInput) (*Message, error)
	GetChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int32) ([]Message, error)

	CreateChat(ctx context.Context, typeID int32, title, avatarURL *string) (*Chat, error)
	AddChatMember(ctx context.Context, chatID, userID uuid.UUID, role string) error
}

type ChatService interface {
	SendMessage(ctx context.Context, input CreateMessageInput) (*Message, error)
	GetHistory(ctx context.Context, chatID, userID uuid.UUID, limit, offset int32) ([]Message, error)
	CreateChat(ctx context.Context, input CreateChatInput) (*Chat, error)
}
