package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/shopspring/decimal"
)

func (s *AccountService) Withdraw(
	ctx context.Context,
	input domain.ServiceWithdrawInput,
) (*domain.WithdrawResult, error) {

	acc, sourceTypeID, sourceUUID, err := s.resolveAccount(ctx, input.SourceType, input.SourceValue)
	if err != nil {
		return nil, fmt.Errorf("source resolution failed: %w", err)
	}

	amount, err := decimal.NewFromString(input.AmountStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidAmountFormat, err)
	}
	if !amount.IsPositive() {
		return nil, domain.ErrInvalidWithdrawAmount
	}

	params := domain.RepoWithdrawParams{
		SourceTypeID:   sourceTypeID,
		SourceID:       sourceUUID,
		AccountID:      acc.ID,
		Amount:         amount,
		IdempotencyKey: input.IdempotencyKey,
	}

	result, err := s.repo.WithdrawTx(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("withdraw transaction failed: %w", err)
	}

	err = s.publisher.PublishWithdrawalCompleted(ctx, domain.WithdrawalCompletedEvent{
		TransactionID: result.TransactionID,
		SourceType:    input.SourceType,
		SourceID:      sourceUUID,
		Amount:        input.AmountStr,
		Currency:      acc.CurrencyCode,
		Timestamp:     time.Now().UTC(),
	})

	if err != nil {
		log.Printf("ERROR: Failed to publish WithdrawalCompleted event for tx %s: %v\n",
			result.TransactionID, err)
	}

	return result, nil
}
