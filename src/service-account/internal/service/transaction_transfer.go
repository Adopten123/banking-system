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

	fromAcc, sourceTypeID, sourceUUID, err := s.resolveAccount(ctx, input.SourceType, input.SourceID)
	if err != nil {
		return nil, fmt.Errorf("sender resolution failed: %w", err)
	}

	toAcc, destinationTypeID, destinationUUID, err := s.resolveAccount(ctx, input.DestinationType, input.DestinationID)
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
	entityValue string,
) (*domain.Account, int32, uuid.UUID, error) {

	switch entityType {
	case "account":
		accUUID, err := uuid.Parse(entityValue)
		if err != nil {
			return nil, 0, uuid.Nil, fmt.Errorf("invalid account public_id format: %w", err)
		}

		acc, err := s.repo.GetByPublicID(ctx, accUUID)
		if err != nil {
			return nil, 0, uuid.Nil, err
		}
		return acc, domain.SourceTypeIDAccount, accUUID, nil

	case "card":
		cardTokenStr, err := s.vault.GetTokenByPan(ctx, entityValue)
		if err != nil {
			return nil, 0, uuid.Nil, fmt.Errorf("card not found in vault: %w", err)
		}
		cardUUID, err := uuid.Parse(cardTokenStr)
		if err != nil {
			return nil, 0, uuid.Nil, fmt.Errorf("vault returned invalid token: %w", err)
		}
		card, err := s.repo.GetCardByID(ctx, cardUUID)
		if err != nil {
			return nil, 0, uuid.Nil, err
		}
		if card.Status != "active" {
			return nil, 0, uuid.Nil, domain.ErrCardBlocked
		}
		acc, err := s.repo.GetAccountInternalByID(ctx, card.AccountID)
		if err != nil {
			return nil, 0, uuid.Nil, fmt.Errorf("failed to get account for card: %w", err)
		}

		return acc, domain.SourceTypeIDCard, cardUUID, nil

	default:
		return nil, 0, uuid.Nil, fmt.Errorf("unsupported entity type: %s", entityType)
	}
}
