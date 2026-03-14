package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/shopspring/decimal"
)

func (s *AccountService) Transfer(
	ctx context.Context,
	input domain.TransferInput,
) (*domain.TransferResult, error) {

	fromAcc, sourceTypeID, sourceUUID, err := s.resolveAccount(ctx, input.SourceType, input.SourceID)
	if err != nil {
		return nil, fmt.Errorf("sender resolution failed: %w", err)
	}

	toAcc, destinationTypeID, destinationUUID, err := s.resolveAccount(ctx, input.DestinationType, input.DestinationID)
	if err != nil {
		return nil, fmt.Errorf("receiver resolution failed: %w", err)
	}

	if fromAcc.ID == toAcc.ID {
		return nil, domain.ErrTransferToSelf
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
		SourceID:          sourceUUID,
		DestinationTypeID: destinationTypeID,
		DestinationID:     destinationUUID,
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

		SourceType: input.SourceType,
		SourceID:   sourceUUID,

		DestinationType: input.DestinationType,
		DestinationID:   destinationUUID,

		SenderAmount:     input.Amount,
		SenderCurrency:   fromAcc.CurrencyCode,
		ReceiverAmount:   receiverAmount.String(),
		ReceiverCurrency: toAcc.CurrencyCode,
		ExchangeRate:     exchangeRate.String(),

		IdempotencyKey: input.IdempotencyKey,
		Timestamp:      time.Now().UTC(),
	}

	err = s.publisher.PublishTransferCreated(ctx, event)
	if err != nil {
		log.Printf("ERROR: Failed to publish TransferCreated event for key %s: %v\n",
			input.IdempotencyKey, err)
	}

	return result, nil
}
