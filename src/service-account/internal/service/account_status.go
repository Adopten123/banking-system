package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
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

	err = s.publisher.PublishAccountStatusChanged(ctx, domain.AccountStatusChangedEvent{
		AccountID: acc.ID,
		OldStatus: acc.StatusID,
		NewStatus: statusBlocked,
		Timestamp: time.Now().UTC(),
	})

	if err != nil {
		log.Printf("ERROR: Failed to publish AccountStatusChanged (Block) event for account %d: %v\n", acc.ID, err)
	}

	return nil
}

func (s *AccountService) CloseAccount(ctx context.Context, publicID uuid.UUID) error {
	const statusClosed = 4

	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	err = s.repo.CloseAccountTx(ctx, acc.ID)
	if err != nil {
		return fmt.Errorf("account didn't close: %w", err)
	}
	err = s.publisher.PublishAccountStatusChanged(ctx, domain.AccountStatusChangedEvent{
		AccountID: acc.ID,
		OldStatus: acc.StatusID,
		NewStatus: statusClosed,
		Timestamp: time.Now().UTC(),
	})

	if err != nil {
		log.Printf("ERROR: Failed to publish AccountStatusChanged (Close) event for account %d: %v\n", acc.ID, err)
	}

	return nil
}
