package service

import (
	"context"

	"github.com/Adopten123/banking-system/service-account/internal/domain"

	"github.com/google/uuid"
)

func (s *AccountService) Transfer(
	ctx context.Context,
	fromPublicID, toPublicID uuid.UUID,
	amount, currency, idempotencyKey, description string,
) error {

	fromAcc, err := s.repo.GetByPublicID(ctx, fromPublicID)
	if err != nil {
		return err
	}

	toAcc, err := s.repo.GetByPublicID(ctx, toPublicID)
	if err != nil {
		return err
	}

	params := domain.TransferParams{
		FromAccountID:  fromAcc.ID,
		ToAccountID:    toAcc.ID,
		AmountStr:      amount,
		CurrencyCode:   currency,
		IdempotencyKey: idempotencyKey,
		Description:    description,
	}
	return s.repo.TransferTx(ctx, params)
}
