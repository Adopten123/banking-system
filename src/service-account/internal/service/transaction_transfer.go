package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

func (s *AccountService) Transfer(
	ctx context.Context,
	input domain.TransferInput,
) (*domain.TransferResult, error) {

	fromAcc, err := s.repo.GetByPublicID(ctx, input.FromPublicID)
	if err != nil {
		return nil, err
	}

	toAcc, err := s.repo.GetByPublicID(ctx, input.ToPublicID)
	if err != nil {
		return nil, err
	}

	params := domain.TransferParams{
		FromAccountID:  fromAcc.ID,
		ToAccountID:    toAcc.ID,
		AmountStr:      input.Amount,
		CurrencyCode:   input.Currency,
		IdempotencyKey: input.IdempotencyKey,
		Description:    input.Description,
	}

	result, err := s.repo.TransferTx(ctx, params)
	if err != nil {
		return nil, err
	}

	event := domain.TransferCreatedEvent{
		TransactionID:  uuid.New(),
		FromAccountID:  fromAcc.ID,
		ToAccountID:    toAcc.ID,
		Amount:         input.Amount,
		Currency:       input.Currency,
		IdempotencyKey: input.IdempotencyKey,
		Timestamp:      time.Now().UTC(),
	}

	err = s.publisher.PublishTransferCreated(ctx, event)
	if err != nil {
		fmt.Printf("ERROR: Failed to publish transfer event for idempotency key %s: %v\n", input.IdempotencyKey, err)
	}

	return result, nil
}
