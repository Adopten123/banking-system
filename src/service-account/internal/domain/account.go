package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID           int64
	PublicID     uuid.UUID
	UserID       uuid.UUID
	TypeID       int32
	StatusID     int32
	CurrencyCode string
	Name         string
	Version      int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AccountInfoResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	TypeID       int32     `json:"type_id"`
	StatusID     int32     `json:"status_id"`
	CurrencyCode string    `json:"currency_code"`
	Name         string    `json:"name"`
	CreatedAt    string    `json:"created_at"`
}

type AccountBalanceResponse struct {
	AccountID uuid.UUID       `json:"account_id"`
	Balance   decimal.Decimal `json:"balance"`
}
