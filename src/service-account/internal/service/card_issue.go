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
		return nil, fmt.Errorf("cannot issue card: account is not active")
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

	expiryDate := time.Date(int(vaultData.ExpiryYear), time.Month(vaultData.ExpiryMonth), 1, 0, 0, 0, 0, time.UTC)

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
		return nil, fmt.Errorf("failed to save card to database: %w", err)
	}
	return card, nil
}
