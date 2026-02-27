package domain

import (
	"github.com/google/uuid"
)

// CreateAccountInput - data for account creating
type CreateAccountInput struct {
	UserID       uuid.UUID
	TypeID       int32
	CurrencyCode string
	Name         string
}
