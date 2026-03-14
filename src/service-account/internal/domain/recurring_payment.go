package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RecurringPayment struct {
	ID                uuid.UUID
	SourceTypeID      int32
	SourceID          uuid.UUID
	DestinationTypeID int32
	DestinationID     uuid.UUID
	Amount            decimal.Decimal
	CurrencyCode      string
	CategoryID        int32
	CronExpression    string
	NextExecutionTime time.Time
	IsActive          bool
	Description       string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CreateRecurringPaymentInput struct {
	SourceType       string
	SourceValue      string
	DestinationType  string
	DestinationValue string
	Amount           string
	CurrencyCode     string
	CategoryID       int32
	CronExpression   string
	Description      string
}

type CreateRecurringPaymentRequest struct {
	SourceType       string `json:"source_type"`
	SourceID         string `json:"source_id"`
	DestinationType  string `json:"destination_type"`
	DestinationID    string `json:"destination_id"`
	Amount           string `json:"amount"`
	CurrencyCode     string `json:"currency_code"`
	CategoryID       int32  `json:"category_id"`
	CronExpression   string `json:"cron_expression"`
	Description      string `json:"description"`
}