package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

func (s *AccountService) IssueCard(ctx context.Context, input domain.IssueCardInput) (*domain.Card, error) {
	const activeStatus = 1

	acc, err := s.repo.GetByPublicID(ctx, input.AccountPublicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if acc.StatusID != activeStatus {
		return nil, domain.ErrAccountInactive
	}

	vaultData, err := s.vault.IssueCard(ctx, domain.IssueCardParams{
		PaymentSystem: input.PaymentSystem,
		IsVirtual:     input.IsVirtual,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to issue secure card token: %w", err)
	}

	tokenUUID, err := uuid.Parse(vaultData.TokenID)
	if err != nil {
		return nil, fmt.Errorf("vault returned invalid token UUID: %w", err)
	}

	nextMonth := time.Date(int(vaultData.ExpiryYear), time.Month(vaultData.ExpiryMonth)+1, 1, 0, 0, 0, 0, time.UTC)
	expiryDate := nextMonth.Add(-time.Second)

	card := &domain.Card{
		ID:        tokenUUID,
		AccountID: acc.ID,
		PANMask:   vaultData.PANMask,
		Expiry:    expiryDate,
		IsVirtual: input.IsVirtual,
		Status:    "active",
		CreatedAt: time.Now().UTC(),
	}

	err = s.repo.CreateCard(ctx, card)
	if err != nil {
		rollbackErr := s.vault.DeleteCardData(ctx, vaultData.TokenID)
		if rollbackErr != nil {
			fmt.Printf("CRITICAL: Failed to rollback vault for token %s: %v\n", vaultData.TokenID, rollbackErr)
		}
		return nil, fmt.Errorf("failed to save card to database: %w", err)
	}

	err = s.publisher.PublishCardIssued(ctx, domain.CardIssuedEvent{
		CardID:    card.ID,
		AccountID: card.AccountID,
		PanMask:   card.PANMask,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		fmt.Printf("ERROR: Failed to publish CardIssued event: %v\n", err)
	}

	return card, nil
}
