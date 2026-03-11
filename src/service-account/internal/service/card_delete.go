package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *AccountService) DeleteCard(ctx context.Context, cardID uuid.UUID) error {
	_, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return err
	}

	err = s.vault.DeleteCardData(ctx, cardID.String())
	if err != nil {
		return fmt.Errorf("failed to delete secure data from vault: %w", err)
	}

	err = s.repo.UpdateCardStatus(ctx, cardID, "deleted")
	if err != nil {
		fmt.Printf("CRITICAL: Vault data deleted but local DB update failed for card %s: %v\n", cardID, err)
		return fmt.Errorf("failed to update local database status to deleted: %w", err)
	}

	return nil
}