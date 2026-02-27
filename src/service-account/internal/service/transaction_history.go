package service

import (
	"context"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

func (s *AccountService) GetAccountTransactions(
	ctx context.Context,
	publicID uuid.UUID,
	params domain.TransactionFilter,
) ([]domain.TransactionHistory, error) {
	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetTransactions(ctx, acc.ID, params)
}
