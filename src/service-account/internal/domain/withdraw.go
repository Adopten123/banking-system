package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WithdrawRequest - body of withdraw funds
type WithdrawRequest struct {
	Amount string `json:"amount"`
}

type WithdrawResponse struct {
	TransactionID uuid.UUID       `json:"transaction_id"`
	NewBalance    decimal.Decimal `json:"new_balance"`
	Currency      string          `json:"currency"`
}
