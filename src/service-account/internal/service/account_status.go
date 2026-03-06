package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

const (
	StatusActive  = 1
	StatusFrozen  = 2
	StatusBlocked = 3
	StatusClosed  = 4
)

func (s *AccountService) changeStatus(ctx context.Context, publicID uuid.UUID, newStatus int32) error {
	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	if acc.StatusID == newStatus {
		return nil
	}

	if acc.StatusID == StatusClosed {
		return fmt.Errorf("cannot change status of a closed account")
	}

	err = s.repo.UpdateStatus(ctx, acc.ID, newStatus)
	if err != nil {
		return fmt.Errorf("failed to update status to %d: %w", newStatus, err)
	}

	err = s.publisher.PublishAccountStatusChanged(ctx, domain.AccountStatusChangedEvent{
		AccountID: acc.ID,
		OldStatus: acc.StatusID,
		NewStatus: newStatus,
		Timestamp: time.Now().UTC(),
	})

	if err != nil {
		log.Printf("ERROR: Failed to publish AccountStatusChanged event for account %d: %v\n", acc.ID, err)
	}

	return nil
}

func (s *AccountService) BlockAccount(ctx context.Context, publicID uuid.UUID) error {
	return s.changeStatus(ctx, publicID, StatusBlocked)
}

func (s *AccountService) FreezeAccount(ctx context.Context, publicID uuid.UUID) error {
	return s.changeStatus(ctx, publicID, StatusFrozen)
}

func (s *AccountService) ActivateAccount(ctx context.Context, publicID uuid.UUID) error {
	return s.changeStatus(ctx, publicID, StatusActive)
}

func (s *AccountService) CloseAccount(ctx context.Context, publicID uuid.UUID) error {
	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	err = s.repo.CloseAccountTx(ctx, acc.ID)
	if err != nil {
		return err
	}

	err = s.publisher.PublishAccountStatusChanged(ctx, domain.AccountStatusChangedEvent{
		AccountID: acc.ID,
		OldStatus: acc.StatusID,
		NewStatus: StatusClosed,
		Timestamp: time.Now().UTC(),
	})

	if err != nil {
		log.Printf("ERROR: Failed to publish AccountStatusChanged (Close) event for account %d: %v\n", acc.ID, err)
	}

	return nil
}
