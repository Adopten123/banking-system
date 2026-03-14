package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/shopspring/decimal"
)

func (s *AccountService) CreateRecurringPayment(ctx context.Context, input domain.CreateRecurringPaymentInput) (uuid.UUID, error) {
	amount, err := decimal.NewFromString(input.Amount)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", domain.ErrInvalidAmountFormat, err)
	}
	if !amount.IsPositive() {
		return uuid.Nil, fmt.Errorf("recurring payment amount must be positive")
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(input.CronExpression)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid cron expression: %w", err)
	}

	nextExecutionTime := schedule.Next(time.Now().UTC())

	fromAcc, sourceTypeID, sourceUUID, err := s.resolveAccount(ctx, input.SourceType, input.SourceValue)
	if err != nil {
		return uuid.Nil, fmt.Errorf("source resolution failed: %w", err)
	}

	currencyCode := input.CurrencyCode
	if currencyCode == "" {
		currencyCode = fromAcc.CurrencyCode
	}

	var destTypeID int32
	var destUUID uuid.UUID
	if input.DestinationType != "" && input.DestinationValue != "" {
		_, destTypeID, destUUID, err = s.resolveAccount(ctx, input.DestinationType, input.DestinationValue)
		if err != nil {
			return uuid.Nil, fmt.Errorf("destination resolution failed: %w", err)
		}
	}

	paymentID := uuid.New()
	payment := domain.RecurringPayment{
		ID:                paymentID,
		SourceTypeID:      sourceTypeID,
		SourceID:          sourceUUID,
		DestinationTypeID: destTypeID,
		DestinationID:     destUUID,
		Amount:            amount,
		CurrencyCode:      currencyCode,
		CategoryID:        input.CategoryID,
		CronExpression:    input.CronExpression,
		NextExecutionTime: nextExecutionTime,
		Description:       input.Description,
	}

	err = s.repo.CreateRecurringPayment(ctx, payment)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to save recurring payment: %w", err)
	}

	return paymentID, nil
}

func (s *AccountService) CancelRecurringPayment(ctx context.Context, paymentID uuid.UUID) error {
	err := s.repo.CancelRecurringPayment(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to cancel recurring payment: %w", err)
	}
	return nil
}