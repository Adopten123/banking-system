package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (r *AccountRepo) WithdrawTx(
	ctx context.Context,
	params domain.RepoWithdrawParams,
) (*domain.WithdrawResult, error) {

	const activeAccountStatus = 1
	const statusCompleted = 1
	const transactionTypeWithdraw = 2

	// Open transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Block account and get actual data
	accInfo, err := qtx.GetAccountForUpdate(ctx, params.AccountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to lock account: %w", err)
	}

	// Check account status
	if accInfo.StatusID.Int32 != activeAccountStatus {
		return nil, domain.ErrAccountInactive
	}

	// Parse balance and limit
	var balance, creditLimit decimal.Decimal
	if val, err := accInfo.Balance.Value(); err == nil && val != nil {
		_ = balance.Scan(val)
	}
	if val, err := accInfo.CreditLimit.Value(); err == nil && val != nil {
		_ = creditLimit.Scan(val)
	}

	// Balance + Credit Limit >= Amount
	availableFunds := balance.Add(creditLimit)
	if availableFunds.LessThan(params.Amount) {
		return nil, domain.ErrInsufficientFunds
	}

	// Subtract money from balance
	var amountNumeric pgtype.Numeric
	_ = amountNumeric.Scan(params.Amount.String())

	err = qtx.SubtractAccountBalance(ctx, SubtractAccountBalanceParams{
		Balance:   amountNumeric,
		AccountID: accInfo.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Make transaction
	txID := uuid.New()
	pgTxID := pgtype.UUID{Bytes: txID, Valid: true}
	var pgSourceID pgtype.UUID
	_ = pgSourceID.Scan(params.SourceID.String())

	_, err = qtx.CreateTransaction(ctx, CreateTransactionParams{
		ID:                pgTxID,
		SourceTypeID:      pgtype.Int4{Int32: params.SourceTypeID, Valid: true},
		SourceID:          pgSourceID,
		DestinationTypeID: pgtype.Int4{Valid: false},
		DestinationID:     pgtype.UUID{Valid: false},
		CategoryID:        pgtype.Int4{Int32: transactionTypeWithdraw, Valid: true},
		StatusID:          pgtype.Int4{Int32: statusCompleted, Valid: true},
		Description:       pgtype.Text{String: "ATM Withdrawal", Valid: true},
		IdempotencyKey:    pgtype.Text{String: params.IdempotencyKey, Valid: true},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrDuplicateTransaction
		}
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Make posting
	negAmount := params.Amount.Neg()
	var negAmountNumeric pgtype.Numeric
	_ = negAmountNumeric.Scan(negAmount.String())

	var exchangeRate pgtype.Numeric
	_ = exchangeRate.Scan("1.0")

	_, err = qtx.CreatePosting(ctx, CreatePostingParams{
		TransactionID: pgTxID,
		AccountID:     pgtype.Int8{Int64: accInfo.ID, Valid: true},
		Amount:        negAmountNumeric,
		CurrencyCode:  pgtype.Text{String: accInfo.CurrencyCode.String, Valid: true},
		ExchangeRate:  exchangeRate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create posting record: %w", err)
	}

	newBalance := balance.Sub(params.Amount)
	// Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return &domain.WithdrawResult{
		TransactionID: txID,
		NewBalance:    newBalance,
		Currency:      accInfo.CurrencyCode.String,
	}, nil
}
