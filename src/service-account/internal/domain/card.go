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

// CardDetails - full data about card
type CardDetails struct {
	PAN         string `json:"pan"`
	CVV         string `json:"cvv"`
	ExpiryMonth int32  `json:"expiry_month"`
	ExpiryYear  int32  `json:"expiry_year"`
}
