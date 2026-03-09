package domain

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID        uuid.UUID
	AccountID int64
	PANMask   string
	Expiry    time.Time
	IsVirtual bool
	Status    string
	CreatedAt time.Time
}

// IssueCardInput - data returned by HTTP-handler
type IssueCardInput struct {
	AccountPublicID uuid.UUID
	PaymentSystem   string
	IsVirtual       bool
}

type IssueCardRequest struct {
	PaymentSystem string `json:"payment_system"`
	IsVirtual     bool   `json:"is_virtual"`
}
