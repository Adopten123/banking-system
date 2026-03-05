package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *AccountService) GetAccountBalance(ctx context.Context, publicID uuid.UUID) (decimal.Decimal, error) {
	balance, err := s.repo.GetBalance(ctx, publicID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}