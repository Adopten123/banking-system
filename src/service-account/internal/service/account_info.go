package service

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

func (s *AccountService) GetAccount(ctx context.Context, publicID uuid.UUID) (*domain.Account, error) {
	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	return acc, nil
}
