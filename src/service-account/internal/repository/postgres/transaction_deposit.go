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

func (r *AccountRepo) DepositTx(
	ctx context.Context,
	params domain.RepoDepositParams,
) (*domain.DepositResult, error) {
	// Open transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	var amount pgtype.Numeric
	if err := amount.Scan(params.AmountStr); err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}

	txUUID := uuid.New()
	pgUUID := pgtype.UUID{
		Bytes: txUUID,
		Valid: true,
	}

	var rate pgtype.Numeric
	if err := rate.Scan("1"); err != nil {
		return nil, fmt.Errorf("failed to scan rate: %w", err)
	}

	// Make transactions
	_, err = qtx.CreateTransaction(ctx, CreateTransactionParams{
		ID:             pgUUID,
		CategoryID:     pgtype.Int4{Int32: 1, Valid: true},
		StatusID:       pgtype.Int4{Int32: 2, Valid: true},
		Description:    pgtype.Text{String: "ATM Deposit", Valid: true},
		IdempotencyKey: pgtype.Text{String: params.IdempotencyKey, Valid: true},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrDuplicateTransaction
		}
		return nil, fmt.Errorf("failed to insert transaction: %w", err)
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
		return nil, fmt.Errorf("failed to insert posting: %w", err)
	}

	// Update balance
	updatedAcc, err := qtx.AddAccountBalance(ctx, AddAccountBalanceParams{
		Balance:   amount,
		AccountID: params.AccountID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	// Commit updates
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit deposit transaction: %w", err)
	}

	newBalance, _ := decimal.NewFromString(updatedAcc.Balance.Int.String())

	return &domain.DepositResult{
		TransactionID: txUUID,
		NewBalance:    newBalance,
	}, nil
}
