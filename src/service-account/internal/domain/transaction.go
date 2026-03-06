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


// TransactionFilter - params of pagination and history filtering
type TransactionFilter struct {
	Limit     int32
	Offset    int32
	StartDate *time.Time
	EndDate   *time.Time
}

type TransactionHistoryResult struct {
	Transactions []TransactionHistory
	TotalCount   int64
}

type TransactionHistoryResponse struct {
	Data       []TransactionHistory 	   `json:"data"`
	Limit      int32                       `json:"limit"`
	Offset     int32                       `json:"offset"`
	TotalCount int64                       `json:"total_count"`
}