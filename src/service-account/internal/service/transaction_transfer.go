package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/shopspring/decimal"
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
		TransactionID:    result.TransactionID,
		FromAccountID:    fromAcc.ID,
		ToAccountID:      toAcc.ID,

		SenderAmount:     transferAmount.String(),
		SenderCurrency:   fromAcc.CurrencyCode,

		ReceiverAmount:   receiverAmount.String(),
		ReceiverCurrency: toAcc.CurrencyCode,

		ExchangeRate:     exchangeRate.String(),

		IdempotencyKey:   input.IdempotencyKey,
		Timestamp:        time.Now().UTC(),
	}

	err = s.publisher.PublishTransferCreated(ctx, event)
	if err != nil {
		fmt.Printf("ERROR: Failed to publish transfer event for idempotency key %s: %v\n", input.IdempotencyKey, err)
	}

	return result, nil
}
