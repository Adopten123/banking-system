package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DomainEvent interface {
	EventName() string
}

type TransferCreatedEvent struct {
	TransactionID    uuid.UUID `json:"transaction_id"`
	FromAccountID    int64     `json:"from_account_id"`
	ToAccountID      int64     `json:"to_account_id"`

	SenderAmount     string    `json:"sender_amount"`
	SenderCurrency   string    `json:"sender_currency"`

	ReceiverAmount   string    `json:"receiver_amount"`
	ReceiverCurrency string    `json:"receiver_currency"`

	ExchangeRate     string    `json:"exchange_rate"`

	IdempotencyKey   string    `json:"idempotency_key"`
	Timestamp        time.Time `json:"timestamp"`
}

func (e TransferCreatedEvent) EventName() string {
	return "TransferCreatedEvent"
}

type AccountCreatedEvent struct {
	AccountID int64     `json:"account_id"`
	PublicID  uuid.UUID `json:"public_id"`
	UserID    uuid.UUID `json:"user_id"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

func (e AccountCreatedEvent) EventName() string {
	return "AccountCreatedEvent"
}

type AccountStatusChangedEvent struct {
	AccountID int64     `json:"account_id"`
	OldStatus int32     `json:"old_status"`
	NewStatus int32     `json:"new_status"`
	Timestamp time.Time `json:"timestamp"`
}

func (e AccountStatusChangedEvent) EventName() string {
	return "AccountStatusChangedEvent"
}

type DepositCompletedEvent struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	AccountID     int64     `json:"account_id"`
	Amount        string    `json:"amount"`
	Currency      string    `json:"currency"`
	Timestamp     time.Time `json:"timestamp"`
}

func (e DepositCompletedEvent) EventName() string {
	return "DepositCompletedEvent"
}

type WithdrawalCompletedEvent struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	AccountID     int64     `json:"account_id"`
	Amount        string    `json:"amount"`
	Currency      string    `json:"currency"`
	Timestamp     time.Time `json:"timestamp"`
}

func (e WithdrawalCompletedEvent) EventName() string {
	return "WithdrawalCompletedEvent"
}

type CreditLimitChangedEvent struct {
	AccountID uuid.UUID `json:"account_id"`
	OldLimit  string    `json:"old_limit"`
	NewLimit  string    `json:"new_limit"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

func (e CreditLimitChangedEvent) EventName() string {
	return "CreditLimitChangedEvent"
}

type EventPublisher interface {
	PublishTransferCreated(ctx context.Context, event TransferCreatedEvent) error
	PublishAccountCreated(ctx context.Context, event AccountCreatedEvent) error
	PublishAccountStatusChanged(ctx context.Context, event AccountStatusChangedEvent) error
	PublishDepositCompleted(ctx context.Context, event DepositCompletedEvent) error
	PublishWithdrawalCompleted(ctx context.Context, event WithdrawalCompletedEvent) error
	PublishCreditLimitChanged(ctx context.Context, event CreditLimitChangedEvent) error
}
