package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	SourceTypeIDAccount int32 = 1
	SourceTypeIDCard    int32 = 2
)

type TransferParams struct {
	SourceTypeID int32
	SourceID     uuid.UUID

	DestinationTypeID int32
	DestinationID     uuid.UUID

	FromAccountID int64
	ToAccountID   int64

	SenderAmountStr   string
	ReceiverAmountStr string

	ExchangeRateStr  string
	CurrencyCode     string
	ReceiverCurrency string

	IdempotencyKey string
	Description    string
}

// TransferInput - data for transfers
type TransferInput struct {
	SourceType string
	SourceID   uuid.UUID

	DestinationType string
	DestinationID   uuid.UUID

	Amount         string
	Currency       string
	IdempotencyKey string
	Description    string
}

type TransferRequest struct {
	SourceType      string `json:"source_type"`
	SourceID        string `json:"source_id"`

	DestinationType string `json:"destination_type"`
	DestinationID   string `json:"destination_id"`

	Amount          string `json:"amount"`
	CurrencyCode    string `json:"currency_code"`
	Description     string `json:"description"`
}

type TransferResponse struct {
	TransactionID uuid.UUID       `json:"transaction_id"`
	NewBalance    decimal.Decimal `json:"sender_new_balance"`
}

type TransferResult struct {
	TransactionID    uuid.UUID
	SenderNewBalance decimal.Decimal
}
