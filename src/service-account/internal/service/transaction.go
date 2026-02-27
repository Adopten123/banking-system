package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
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

	// Checking business rule: account must be active
	if acc.StatusID != 1 {
		return domain.ErrAccountInactive
	}

	// Calling repo layer
	err = s.repo.Deposit(ctx, domain.RepoDepositParams{
		AccountID: acc.ID,
		AmountStr: input.AmountStr,
		CurrencyCode: acc.CurrencyCode,
		IdempotencyKey: input.IdempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("failed to process deposit in repository: %w", err)
	}
	return nil
}
