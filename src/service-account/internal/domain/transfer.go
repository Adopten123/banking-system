package domain

import "github.com/google/uuid"

type TransferParams struct {
	FromAccountID  int64
	ToAccountID    int64
	AmountStr      string
	CurrencyCode   string
	IdempotencyKey string
	Description    string
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
