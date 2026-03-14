package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WithdrawRequest - body of withdraw funds
type WithdrawRequest struct {
	SourceType string `json:"source_type"`
	SourceID   string `json:"source_id"`
	Amount     string `json:"amount"`
}

type WithdrawResult struct {
	TransactionID uuid.UUID       `json:"transaction_id"`
	NewBalance    decimal.Decimal `json:"new_balance"`
	Currency      string          `json:"currency"`
}

type ServiceWithdrawInput struct {
	SourceType     string
	SourceValue    string
	AmountStr      string
	IdempotencyKey string
}

type RepoWithdrawParams struct {
	SourceTypeID   int32
	SourceID       uuid.UUID
	AccountID      int64
	Amount         decimal.Decimal
	IdempotencyKey string
}

type WithdrawResponse struct {
	TransactionID uuid.UUID       `json:"transaction_id"`
	NewBalance    decimal.Decimal `json:"new_balance"`
	Currency      string          `json:"currency"`
}
