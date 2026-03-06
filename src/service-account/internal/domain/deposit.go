package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ServiceDepositInput - data for Deposit (Service Layer)
type ServiceDepositInput struct {
	AmountStr      string
	IdempotencyKey string
}

// RepoDepositParams - data for Deposit (Service Repository)
type RepoDepositParams struct {
	AccountID      int64
	AmountStr      string
	CurrencyCode   string
	IdempotencyKey string
}

// DepositResult - data info for JSON Repo - Service - Handler
type DepositResult struct {
	TransactionID uuid.UUID
	NewBalance    decimal.Decimal
}

// DepositRequest - waiting sum for deposit
type DepositRequest struct {
	Amount string `json:"amount"`
}

type DepositResponse struct {
	TransactionID uuid.UUID       `json:"transaction_id"`
	NewBalance    decimal.Decimal `json:"new_balance"`
}
