package domain

import (
	"context"

	"github.com/google/uuid"
)

type TransferCreatedEvent struct {
	MessageID      int64     `json:"message_id"`
	ChatID         uuid.UUID `json:"chat_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	Amount         string    `json:"amount"`
	Currency       string    `json:"currency"`
	IdempotencyKey string    `json:"idempotency_key"`
}

type TransferPublisher interface {
	PublishTransferCreated(ctx context.Context, event TransferCreatedEvent) error
}
