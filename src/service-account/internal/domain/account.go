package domain

import (
	"time"

	"github.com/google/uuid"
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
