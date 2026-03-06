package service

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *AccountService) Withdraw(
	ctx context.Context,
	publicID uuid.UUID,
	amount decimal.Decimal,
	idempotencyKey string,
) (*domain.WithdrawResponse, error) {
	result, err := s.repo.WithdrawTx(ctx, publicID, amount, idempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("withdraw transaction failed: %w", err)
	}
	// TODO: RabbitMQ + Kafka
	return result, nil
}
