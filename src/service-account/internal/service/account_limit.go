package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *AccountService) UpdateCreditLimit(ctx context.Context, publicID uuid.UUID, limitStr string) error {
	newLimit, err := decimal.NewFromString(limitStr)
	if err != nil || newLimit.IsNegative() {
		return fmt.Errorf("invalid credit limit: must be a non-negative number")
	}

	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	if acc.StatusID != 1 {
		return domain.ErrAccountInactive
	}

	oldLimitStr, err := s.repo.GetCreditLimit(ctx, acc.ID)
	if err != nil {
		return fmt.Errorf("failed to get old credit limit: %w", err)
	}

	oldLimit, _ := decimal.NewFromString(oldLimitStr)
	if oldLimit.Equal(newLimit) {
		return nil
	}

	err = s.repo.UpdateCreditLimit(ctx, acc.ID, newLimit.String())
	if err != nil {
		return fmt.Errorf("failed to update credit limit: %w", err)
	}

	event := domain.CreditLimitChangedEvent{
		AccountID: publicID,
		OldLimit:  oldLimitStr,
		NewLimit:  newLimit.String(),
		Currency:  acc.CurrencyCode,
		Timestamp: time.Now().UTC(),
	}

	err = s.publisher.PublishCreditLimitChanged(ctx, event)
	if err != nil {
		fmt.Printf("ERROR: Failed to publish CreditLimitChanged event for account %s: %v\n", publicID, err)
	}

	return nil
}
