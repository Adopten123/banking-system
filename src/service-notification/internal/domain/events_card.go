package domain

import (
	"time"

	"github.com/google/uuid"
)

type CardIssuedEvent struct {
	CardID    uuid.UUID `json:"card_id"`
	AccountID int64     `json:"account_id"`
	PanMask   string    `json:"pan_mask"`
	Timestamp time.Time `json:"timestamp"`
}

type CardStatusChangedEvent struct {
	CardID    uuid.UUID `json:"card_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	Timestamp time.Time `json:"timestamp"`
}