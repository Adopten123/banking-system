package service

import (
	"context"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
)

func (s *AccountService) Transfer(
	ctx context.Context,
	input domain.TransferInput,
) error {

	fromAcc, err := s.repo.GetByPublicID(ctx, input.FromPublicID)
	if err != nil {
		return err
	}

	toAcc, err := s.repo.GetByPublicID(ctx, input.ToPublicID)
	if err != nil {
		return err
	}

	params := domain.TransferParams{
		FromAccountID:  fromAcc.ID,
		ToAccountID:    toAcc.ID,
		AmountStr:      input.Amount,
		CurrencyCode:   input.Currency,
		IdempotencyKey: input.IdempotencyKey,
		Description:    input.Description,
	}
	return s.repo.TransferTx(ctx, params)
}
