package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

func (s *AccountService) UpdateCardStatus(ctx context.Context, cardID uuid.UUID, newStatus string) error {
	// 1. Validate data
	if newStatus != "active" && newStatus != "blocked" {
		return domain.ErrInvalidCardStatus
	}

	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return err
	}

	if card.Status == newStatus {
		return nil
	}

	// 2. Update status in vault
	err = s.repo.UpdateCardStatus(ctx, cardID, newStatus)
	if err != nil {
		fmt.Printf("CRITICAL: Vault updated but local DB failed for card %s: %v\n", cardID, err)
		return fmt.Errorf("failed to update local database: %w", err)
	}

	// 3. Update status in card service
	err = s.publisher.PublishCardStatusChanged(ctx, domain.CardStatusChangedEvent{
		CardID:    cardID,
		OldStatus: card.Status,
		NewStatus: newStatus,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		fmt.Printf("ERROR: Failed to publish CardStatusChanged event: %v\n", err)
	}

	return nil
}
