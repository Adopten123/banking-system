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
	publicID uuid.UUID,
	amount decimal.Decimal,
	idempotencyKey string,
) (*domain.WithdrawResponse, error) {

	const activeAccountStatus = 1
	const statusCompleted = 1
	const transactionTypeWithdraw = 3

	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(publicID.String()); err != nil {
		return nil, fmt.Errorf("invalid public_id: %w", err)
	}

	// Open transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Block account and get actual data
	accInfo, err := qtx.GetAccountForWithdrawUpdate(ctx, pgUUID)
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
	balance, err := decimal.NewFromString(accInfo.AbBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to parse balance: %w", err)
	}
	creditLimit, err := decimal.NewFromString(accInfo.AbCreditLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credit limit: %w", err)
	}

	// Balance + Credit Limit >= Amount
	availableFunds := balance.Add(creditLimit)
	if availableFunds.LessThan(amount) {
		return nil, domain.ErrInsufficientFunds
	}

	// Subtract money from balance
	var amountNumeric pgtype.Numeric
	if err := amountNumeric.Scan(amount.String()); err != nil {
		return nil, fmt.Errorf("failed to scan amount: %w", err)
	}

	err = qtx.SubtractAccountBalance(ctx, SubtractAccountBalanceParams{
		Balance:   amountNumeric,
		AccountID: accInfo.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Make transaction
	txID := uuid.New()
	var pgTxID pgtype.UUID
	if err := pgTxID.Scan(txID.String()); err != nil {
		return nil, fmt.Errorf("failed to scan tx uuid: %w", err)
	}

	_, err = qtx.CreateTransaction(ctx, CreateTransactionParams{
		ID:              pgTxID,
		CategoryID:      pgtype.Int4{Int32: transactionTypeWithdraw, Valid: true},
		StatusID:        pgtype.Int4{Int32: statusCompleted, Valid: true},
		Description:     pgtype.Text{String: "ATM Withdrawal", Valid: true},
		ExternalDetails: []byte("{}"),
		IdempotencyKey:  pgtype.Text{String: idempotencyKey, Valid: true},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrDuplicateTransaction
		}
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Make posting
	negAmount := amount.Neg()
	var negAmountNumeric pgtype.Numeric
	if err := negAmountNumeric.Scan(negAmount.String()); err != nil {
		return nil, fmt.Errorf("failed to scan negative amount: %w", err)
	}

	var exchangeRate pgtype.Numeric
	if err := exchangeRate.Scan("1.0"); err != nil {
		return nil, fmt.Errorf("failed to scan exchange rate: %w", err)
	}

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

	newBalance := balance.Sub(amount)
	// Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return &domain.WithdrawResponse{
		TransactionID: txID,
		NewBalance:    newBalance,
		Currency:      accInfo.CurrencyCode.String,
	}, nil
}
