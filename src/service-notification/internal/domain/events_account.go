package domain

import (
	"time"

	"github.com/google/uuid"
)

type AccountCreatedEvent struct {
	AccountID int64     `json:"account_id"`
	PublicID  uuid.UUID `json:"public_id"`
	UserID    uuid.UUID `json:"user_id"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

type AccountStatusChangedEvent struct {
	AccountID int64     `json:"account_id"`
	OldStatus int32     `json:"old_status"`
	NewStatus int32     `json:"new_status"`
	Timestamp time.Time `json:"timestamp"`
}

type CreditLimitChangedEvent struct {
	AccountID uuid.UUID `json:"account_id"`
	OldLimit  string    `json:"old_limit"`
	NewLimit  string    `json:"new_limit"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}