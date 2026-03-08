package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransferParams struct {
	FromAccountID     int64
	ToAccountID       int64
	SenderAmountStr   string
	ReceiverAmountStr string
	ExchangeRateStr   string
	CurrencyCode      string
	ReceiverCurrency  string
	IdempotencyKey    string
	Description       string
}

// TransferInput - data for transfers
type TransferInput struct {
	FromPublicID   uuid.UUID
	ToPublicID     uuid.UUID
	Amount         string
	Currency       string
	IdempotencyKey string
	Description    string
}

type TransferRequest struct {
	ToAccountID  string `json:"to_account_id"`
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
	Description  string `json:"description"`
}

type TransferResponse struct {
	TransactionID uuid.UUID       `json:"transaction_id"`
	NewBalance    decimal.Decimal `json:"sender_new_balance"`
}

type TransferResult struct {
	TransactionID    uuid.UUID
	SenderNewBalance decimal.Decimal
}
