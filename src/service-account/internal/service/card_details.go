package service

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

// GetCardDetails возвращает сырые реквизиты карты, запрашивая их у Vault
func (s *AccountService) GetCardDetails(ctx context.Context, cardID uuid.UUID) (*domain.CardDetails, error) {
	// 1. Try to find card in account_db (Go)
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	// 2. Checking status
	if card.Status != "active" {
		return nil, domain.ErrCardBlocked
	}

	// 3. Getting card details in vault
	details, err := s.vault.GetCardDetails(ctx, card.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sensitive card details: %w", err)
	}

	return details, nil
}
