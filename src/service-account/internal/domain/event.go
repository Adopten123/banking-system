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

type AccountCreatedEvent struct {
	AccountID int64     `json:"account_id"`
	PublicID  uuid.UUID `json:"public_id"`
	UserID    uuid.UUID `json:"user_id"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

type AccountStatusChangedEvent struct {
	AccountID int64     `json:"account_id"`
	OldStatus int32     `json:"old_status"`
	NewStatus int32     `json:"new_status"`
	Timestamp time.Time `json:"timestamp"`
}

type DepositCompletedEvent struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	AccountID     int64     `json:"account_id"`
	Amount        string    `json:"amount"`
	Currency      string    `json:"currency"`
	Timestamp     time.Time `json:"timestamp"`
}

type EventPublisher interface {
	PublishTransferCreated(ctx context.Context, event TransferCreatedEvent) error
	PublishAccountCreated(ctx context.Context, event AccountCreatedEvent) error
	PublishAccountStatusChanged(ctx context.Context, event AccountStatusChangedEvent) error
	PublishDepositCompleted(ctx context.Context, event DepositCompletedEvent) error
}
