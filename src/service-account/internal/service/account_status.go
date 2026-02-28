package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *AccountService) BlockAccount(ctx context.Context, publicID uuid.UUID) error {
	const statusBlocked = 3

	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	err = s.repo.UpdateStatus(ctx, acc.ID, statusBlocked)
	if err != nil {
		return fmt.Errorf("account didn`t block: %w", err)
	}
	return nil
}

func (s *AccountService) CloseAccount(ctx context.Context, publicID uuid.UUID) error {
	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	err = s.repo.CloseAccountTx(ctx, acc.ID)
	if err != nil {
		return fmt.Errorf("account didn't close: %w", err)
	}
	return nil
}
