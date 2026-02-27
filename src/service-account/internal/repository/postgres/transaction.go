package postgres

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *AccountRepo) Deposit(
	ctx context.Context,
	params domain.RepoDepositParams,
) error {
	// Open transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	var amount pgtype.Numeric
	if err := amount.Scan(params.AmountStr); err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	pgUUID := pgtype.UUID{
		Bytes: uuid.New(),
		Valid: true,
	}

	var rate pgtype.Numeric
	_ = rate.Scan("1")

	// Make transactions
	_, err = qtx.CreateTransaction(ctx, CreateTransactionParams{
		ID:             pgUUID,
		CategoryID:     pgtype.Int4{Int32: 1, Valid: true},
		StatusID:       pgtype.Int4{Int32: 2, Valid: true},
		Description:    pgtype.Text{String: "Пополнение счета", Valid: true},
		IdempotencyKey: pgtype.Text{String: params.IdempotencyKey, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	// Make posting
	_, err = qtx.CreatePosting(ctx, CreatePostingParams{
		TransactionID: pgUUID,
		AccountID:     pgtype.Int8{Int64: params.AccountID, Valid: true},
		Amount:        amount,
		CurrencyCode:  pgtype.Text{String: params.CurrencyCode, Valid: true},
		ExchangeRate:  rate,
	})
	if err != nil {
		return fmt.Errorf("failed to insert posting: %w", err)
	}

	// Update balance
	_, err = qtx.AddAccountBalance(ctx, AddAccountBalanceParams{
		Balance:   amount,
		AccountID: params.AccountID,
	})
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	// Commit updates
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit deposit transaction: %w", err)
	}

	return nil
}
