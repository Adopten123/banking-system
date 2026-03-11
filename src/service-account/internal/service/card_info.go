package service

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

func (s *AccountService) GetAccountCards(ctx context.Context, accountPublicID uuid.UUID) ([]*domain.Card, error) {
	acc, err := s.repo.GetByPublicID(ctx, accountPublicID)
	if err != nil {
		return nil, err
	}

	cards, err := s.repo.GetCardsByAccountID(ctx, acc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards for account: %w", err)
	}

	return cards, nil
}