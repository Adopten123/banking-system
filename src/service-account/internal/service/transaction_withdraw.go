package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *AccountService) Withdraw(
	ctx context.Context,
	publicID uuid.UUID,
	amount decimal.Decimal,
	idempotencyKey string,
) (*domain.WithdrawResult, error) {
	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account info: %w", err)
	}

	if acc.StatusID != 1 {
		return nil, domain.ErrAccountInactive
	}

	result, err := s.repo.WithdrawTx(ctx, publicID, amount, idempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("withdraw transaction failed: %w", err)
	}

	err = s.publisher.PublishWithdrawalCompleted(ctx, domain.WithdrawalCompletedEvent{
		TransactionID: result.TransactionID,
		AccountID:     acc.ID,
		Amount:        amount.String(),
		Currency:      acc.CurrencyCode,
		Timestamp:     time.Now().UTC(),
	})
	if err != nil {
		fmt.Printf("ERROR: Failed to publish Withdrawal event for tx %s: %v\n", result.TransactionID, err)
	}

	return result, nil
}
