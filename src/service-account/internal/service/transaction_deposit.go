package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/shopspring/decimal"
)

func (s *AccountService) Deposit(
	ctx context.Context,
	input domain.ServiceDepositInput,
) (*domain.DepositResult, error) {

	acc, destinationTypeID, destinationUUID, err := s.resolveAccount(ctx, input.DestinationType, input.DestinationValue)
	if err != nil {
		return nil, fmt.Errorf("destination resolution failed: %w", err)
	}
	amount, err := decimal.NewFromString(input.AmountStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidAmountFormat, err)
	}
	if !amount.IsPositive() {
		return nil, domain.ErrInvalidDepositAmount
	}

	result, err := s.repo.DepositTx(ctx, domain.RepoDepositParams{
		DestinationTypeID: destinationTypeID,
		DestinationID:     destinationUUID,
		AccountID:         acc.ID,
		AmountStr:         input.AmountStr,
		CurrencyCode:      acc.CurrencyCode,
		IdempotencyKey:    input.IdempotencyKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process deposit in repository: %w", err)
	}

	err = s.publisher.PublishDepositCompleted(ctx, domain.DepositCompletedEvent{
		TransactionID:   result.TransactionID,
		DestinationType: input.DestinationType,
		DestinationID:   destinationUUID,
		Amount:          input.AmountStr,
		Currency:        acc.CurrencyCode,
		Timestamp:       time.Now().UTC(),
	})

	if err != nil {
		log.Printf("ERROR: Failed to publish DepositCompleted event for %s %s: %v\n",
			input.DestinationType, destinationUUID, err)
	}

	return result, nil
}
