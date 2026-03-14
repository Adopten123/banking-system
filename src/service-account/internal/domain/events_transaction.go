package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransferCreatedEvent struct {
	TransactionID uuid.UUID `json:"transaction_id"`

	SourceType      string    `json:"source_type"`
	SourceID        uuid.UUID `json:"source_id"`
	DestinationType string    `json:"destination_type"`
	DestinationID   uuid.UUID `json:"destination_id"`

	SenderAmount     string `json:"sender_amount"`
	SenderCurrency   string `json:"sender_currency"`
	ReceiverAmount   string `json:"receiver_amount"`
	ReceiverCurrency string `json:"receiver_currency"`
	ExchangeRate     string `json:"exchange_rate"`

	IdempotencyKey string    `json:"idempotency_key"`
	Timestamp      time.Time `json:"timestamp"`
}

func (e TransferCreatedEvent) EventName() string { return "TransferCreatedEvent" }

type DepositCompletedEvent struct {
	TransactionID uuid.UUID `json:"transaction_id"`

	DestinationType string    `json:"destination_type"`
	DestinationID   uuid.UUID `json:"destination_id"`

	Amount    string    `json:"amount"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

func (e DepositCompletedEvent) EventName() string { return "DepositCompletedEvent" }

type WithdrawalCompletedEvent struct {
	TransactionID uuid.UUID `json:"transaction_id"`

	SourceType string    `json:"source_type"`
	SourceID   uuid.UUID `json:"source_id"`

	Amount    string    `json:"amount"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

func (e WithdrawalCompletedEvent) EventName() string { return "WithdrawalCompletedEvent" }
