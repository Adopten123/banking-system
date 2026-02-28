package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TransferCreatedEvent struct {
	TransactionID  uuid.UUID `json:"transaction_id"`
	FromAccountID  int64     `json:"from_account_id"`
	ToAccountID    int64     `json:"to_account_id"`
	Amount         string    `json:"amount"`
	Currency       string    `json:"currency"`
	IdempotencyKey string    `json:"idempotency_key"`
	Timestamp      time.Time `json:"timestamp"`
}

type EventPublisher interface {
	PublishTransferCreated(ctx context.Context, event TransferCreatedEvent) error
}
