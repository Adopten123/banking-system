package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (r *AccountRepo) TransferTx(
	ctx context.Context,
	params domain.TransferParams,
) (*domain.TransferResult, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Deadlock defense
	account1ID := params.FromAccountID
	account2ID := params.ToAccountID
	if account1ID > account2ID {
		account1ID, account2ID = account2ID, account1ID
	}

	// Blocking accounts for update
	acc1, err := qtx.GetAccountForUpdate(ctx, account1ID)
	if err != nil {
		return nil, fmt.Errorf("failed to lock account 1: %w", err)
	}

	acc2, err := qtx.GetAccountForUpdate(ctx, account2ID)
	if err != nil {
		return nil, fmt.Errorf("failed to lock account 2: %w", err)
	}

	// Business Checks. Verifying Statuses and Currencies
	var senderAcc, receiverAcc GetAccountForUpdateRow

	if params.FromAccountID == account1ID {
		senderAcc, receiverAcc = acc1, acc2
	} else {
		senderAcc, receiverAcc = acc2, acc1
	}
	if senderAcc.StatusID.Int32 != 1 || receiverAcc.StatusID.Int32 != 1 {
		return nil, domain.ErrAccountInactive
	}

	//if senderAcc.CurrencyCode.String != receiverAcc.CurrencyCode.String {
	//	return nil, errors.New("cross-currency transfers are not supported yet")
	//}

	senderAmount, err := decimal.NewFromString(params.SenderAmountStr)
	if err != nil {
		return nil, fmt.Errorf("invalid sender amount: %w", err)
	}

	receiverAmount, err := decimal.NewFromString(params.ReceiverAmountStr)
	if err != nil {
		return nil, fmt.Errorf("invalid receiver amount: %w", err)
	}

	var senderBalance, senderCreditLimit decimal.Decimal

	if val, err := senderAcc.Balance.Value(); err == nil && val != nil {
		_ = senderBalance.Scan(val)
	}

	if val, err := senderAcc.CreditLimit.Value(); err == nil && val != nil {
		_ = senderCreditLimit.Scan(val)
	}

	if senderBalance.Add(senderCreditLimit).LessThan(senderAmount) {
		return nil, domain.ErrInsufficientFunds
	}

	var amountPositive pgtype.Numeric
	amountPositive.Scan(receiverAmount.String())

	var amountNegative pgtype.Numeric
	amountNegative.Scan(senderAmount.Neg().String())

	var exRate pgtype.Numeric
	exRate.Scan(params.ExchangeRateStr)

	txUUID := uuid.New()
	pgTxID := pgtype.UUID{Bytes: txUUID, Valid: true}

	var pgSourceID pgtype.UUID
	_ = pgSourceID.Scan(params.SourceID.String())

	var pgDestinationID pgtype.UUID
	_ = pgDestinationID.Scan(params.DestinationID.String())

	_, err = qtx.CreateTransaction(ctx, CreateTransactionParams{
		ID:                pgTxID,
		SourceTypeID:      pgtype.Int4{Int32: params.SourceTypeID, Valid: true},
		SourceID:          pgSourceID,
		DestinationTypeID: pgtype.Int4{Int32: params.DestinationTypeID, Valid: true},
		DestinationID:     pgDestinationID,
		IdempotencyKey:    pgtype.Text{String: params.IdempotencyKey, Valid: true},
		CategoryID:        pgtype.Int4{Int32: 3, Valid: true},
		StatusID:          pgtype.Int4{Int32: 2, Valid: true},
		Description:       pgtype.Text{String: params.Description, Valid: true},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrDuplicateTransaction
		}
		return nil, fmt.Errorf("failed to insert transaction: %w", err)
	}

	_, err = qtx.CreatePosting(ctx, CreatePostingParams{
		TransactionID: pgTxID,
		AccountID:     pgtype.Int8{Int64: params.FromAccountID, Valid: true},
		Amount:        amountNegative,
		CurrencyCode:  pgtype.Text{String: params.CurrencyCode, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create posting for sender: %w", err)
	}

	senderUpdatedRow, err := qtx.AddAccountBalance(ctx, AddAccountBalanceParams{
		Balance:   amountNegative,
		AccountID: params.FromAccountID,
	})
	if err != nil {
		return nil, fmt.Errorf("insufficient funds or balance update failed: %w", err)
	}

	_, err = qtx.CreatePosting(ctx, CreatePostingParams{
		TransactionID: pgTxID,
		AccountID:     pgtype.Int8{Int64: params.ToAccountID, Valid: true},
		Amount:        amountPositive,
		CurrencyCode:  pgtype.Text{String: params.CurrencyCode, Valid: true},
		ExchangeRate:  exRate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create posting for receiver: %w", err)
	}

	_, err = qtx.AddAccountBalance(ctx, AddAccountBalanceParams{
		Balance:   amountPositive,
		AccountID: params.ToAccountID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add balance to receiver: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	var senderNewBalance decimal.Decimal
	if val, err := senderUpdatedRow.Balance.Value(); err == nil && val != nil {
		_ = senderNewBalance.Scan(val)
	}

	return &domain.TransferResult{
		TransactionID:    txUUID,
		SenderNewBalance: senderNewBalance,
	}, nil
}
