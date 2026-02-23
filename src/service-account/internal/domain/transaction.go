package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionHistory struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	CategoryID    int32     `json:"category_id"`
	StatusID      int32     `json:"status_id"`
	Description   string    `json:"description"`
	Amount        string    `json:"amount"`
	CurrencyCode  string    `json:"currency_code"`
	CreatedAt     time.Time `json:"created_at"`
}
