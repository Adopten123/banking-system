package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *AccountService) Transfer(
	ctx context.Context,
	input domain.TransferInput,
) (*domain.TransferResult, error) {

	fromAcc, sourceTypeID, err := s.resolveAccount(ctx, input.SourceType, input.SourceID)
	if err != nil {
		return nil, fmt.Errorf("sender resolution failed: %w", err)
	}

	toAcc, destinationTypeID, err := s.resolveAccount(ctx, input.DestinationType, input.DestinationID)
	if err != nil {
		return nil, fmt.Errorf("receiver resolution failed: %w", err)
	}

	transferAmount, err := decimal.NewFromString(input.Amount)
	if err != nil {
		return nil, err
	}

	receiverAmount := transferAmount
	exchangeRate := decimal.NewFromInt(1)

	if fromAcc.CurrencyCode != toAcc.CurrencyCode {
		exchangeRate, err = s.exchanger.GetRate(ctx, fromAcc.CurrencyCode, toAcc.CurrencyCode)
		if err != nil {
			return nil, fmt.Errorf("failed to get exchange rate: %w", err)
		}
		receiverAmount = transferAmount.Mul(exchangeRate).RoundBank(2)
	}

	params := domain.TransferParams{
		SourceTypeID:      sourceTypeID,
		SourceID:          input.SourceID,
		DestinationTypeID: destinationTypeID,
		DestinationID:     input.DestinationID,
		FromAccountID:     fromAcc.ID,
		ToAccountID:       toAcc.ID,
		SenderAmountStr:   transferAmount.String(),
		ReceiverAmountStr: receiverAmount.String(),
		ExchangeRateStr:   exchangeRate.String(),
		CurrencyCode:      fromAcc.CurrencyCode,
		ReceiverCurrency:  toAcc.CurrencyCode,
		IdempotencyKey:    input.IdempotencyKey,
		Description:       input.Description,
	}

	result, err := s.repo.TransferTx(ctx, params)
	if err != nil {
		return nil, err
	}

	event := domain.TransferCreatedEvent{
		TransactionID: result.TransactionID,
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,

		SenderAmount:   transferAmount.String(),
		SenderCurrency: fromAcc.CurrencyCode,

		ReceiverAmount:   receiverAmount.String(),
		ReceiverCurrency: toAcc.CurrencyCode,

		ExchangeRate: exchangeRate.String(),

		IdempotencyKey: input.IdempotencyKey,
		Timestamp:      time.Now().UTC(),
	}

	err = s.publisher.PublishTransferCreated(ctx, event)
	if err != nil {
		fmt.Printf("ERROR: Failed to publish transfer event for idempotency key %s: %v\n", input.IdempotencyKey, err)
	}

	return result, nil
}

// resolveAccount - find account by source
func (s *AccountService) resolveAccount(
	ctx context.Context,
	entityType string,
	entityID uuid.UUID,
) (*domain.Account, int32, error) {

	switch entityType {
	case "account":
		acc, err := s.repo.GetByPublicID(ctx, entityID)
		if err != nil {
			return nil, 0, err
		}
		return acc, domain.SourceTypeIDAccount, nil

	case "card":
		card, err := s.repo.GetCardByID(ctx, entityID)
		if err != nil {
			return nil, 0, err
		}

		if card.Status != "active" {
			return nil, 0, domain.ErrCardBlocked
		}

		acc, err := s.repo.GetAccountInternalByID(ctx, card.AccountID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get account for card: %w", err)
		}
		return acc, domain.SourceTypeIDCard, nil

	default:
		return nil, 0, fmt.Errorf("unsupported entity type: %s", entityType)
	}
}
