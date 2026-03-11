package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *AccountService) UpdateCardStatus(ctx context.Context, cardID uuid.UUID, newStatus string) error {
	// 1. Validate data
	if newStatus != "active" && newStatus != "blocked" {
		return fmt.Errorf("invalid status: %s", newStatus)
	}

	// 2. Update status in vault
	err := s.vault.UpdateCardStatus(ctx, cardID.String(), newStatus)
	if err != nil {
		return fmt.Errorf("failed to sync status with vault: %w", err)
	}

	// 3. Update status in card service
	err = s.repo.UpdateCardStatus(ctx, cardID, newStatus)
	if err != nil {
		fmt.Printf("CRITICAL: Vault updated but local DB failed for card %s: %v\n", cardID, err)
		return fmt.Errorf("failed to update local database: %w", err)
	}

	return nil
}
