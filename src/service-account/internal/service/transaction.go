package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *AccountService) Deposit(
	ctx context.Context,
	publicID uuid.UUID,
	input domain.ServiceDepositInput,
) error {
	// Getting account by id
	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return err
		}
		return fmt.Errorf("failed to fetch account for deposit: %w", err)
	}

	// Checking business rule:
	// Account must be active
	if acc.StatusID != 1 {
		return domain.ErrAccountInactive
	}
	// Amount > 0
	amount, err := decimal.NewFromString(input.AmountStr)
	if err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}
	if !amount.IsPositive() {
		return domain.ErrInvalidDepositAmount
	}

	// Calling repo layer
	err = s.repo.Deposit(ctx, domain.RepoDepositParams{
		AccountID:      acc.ID,
		AmountStr:      input.AmountStr,
		CurrencyCode:   acc.CurrencyCode,
		IdempotencyKey: input.IdempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("failed to process deposit in repository: %w", err)
	}

	err = s.publisher.PublishDepositCompleted(ctx, domain.DepositCompletedEvent{
		TransactionID: uuid.New(),
		AccountID:     acc.ID,
		Amount:        input.AmountStr,
		Currency:      acc.CurrencyCode,
		Timestamp:     time.Now().UTC(),
	})

	if err != nil {
		log.Printf("ERROR: Failed to publish DepositCompleted event for account %d: %v\n", acc.ID, err)
	}

	return nil
}
