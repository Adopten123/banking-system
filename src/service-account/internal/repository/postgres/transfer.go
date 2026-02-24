package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type TransferTxParams struct {
	FromAccountID  int64
	ToAccountID    int64
	AmountStr      string
	CurrencyCode   string
	IdempotencyKey string
	Description    string
}

func (r *AccountRepo) TransferTx(ctx context.Context, arg TransferTxParams) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Deadlock defense
	account1ID := arg.FromAccountID
	account2ID := arg.ToAccountID
	if account1ID > account2ID {
		account1ID, account2ID = account2ID, account1ID
	}

	// Blocking accounts for update
	acc1, err := qtx.GetAccountForUpdate(ctx, account1ID)
	if err != nil {
		return fmt.Errorf("failed to lock account 1: %w", err)
	}

	acc2, err := qtx.GetAccountForUpdate(ctx, account2ID)
	if err != nil {
		return fmt.Errorf("failed to lock account 2: %w", err)
	}

	// Business Checks. Verifying Statuses and Currencies
	var senderStatus int32
	var senderCurrency, receiverCurrency string

	if arg.FromAccountID == account1ID {
		senderStatus = acc1.StatusID.Int32
		senderCurrency = acc1.CurrencyCode.String
		receiverCurrency = acc2.CurrencyCode.String
	} else {
		senderStatus = acc2.StatusID.Int32
		senderCurrency = acc2.CurrencyCode.String
		receiverCurrency = acc1.CurrencyCode.String
	}

	if senderStatus != 1 {
		return errors.New("sender account is blocked or inactive")
	}
	if senderCurrency != receiverCurrency {
		return errors.New("cross-currency transfers are not supported yet")
	}

	// Prepare Sum
	var amountPositive, amountNegative pgtype.Numeric
	amountPositive.Scan(arg.AmountStr)
	amountNegative.Scan("-" + arg.AmountStr)

	// Making transaction
	//(category_id = 3 - transfer, status_id = 2 - posted)
	txID, err := qtx.CreateTransaction(ctx, CreateTransactionParams{
		IdempotencyKey: pgtype.Text{String: arg.IdempotencyKey, Valid: true},
		CategoryID:     pgtype.Int4{Int32: 3, Valid: true},
		StatusID:       pgtype.Int4{Int32: 2, Valid: true},
		Description:    pgtype.Text{String: arg.Description, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Posting and writing off the Sender's balance
	_, err = qtx.CreatePosting(ctx, CreatePostingParams{
		TransactionID: txID.ID,
		AccountID:     pgtype.Int8{Int64: arg.FromAccountID, Valid: true},
		Amount:        amountNegative,
		CurrencyCode:  pgtype.Text{String: arg.CurrencyCode, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create posting for sender: %w", err)
	}

	_, err = qtx.AddAccountBalance(ctx, AddAccountBalanceParams{
		Balance:   amountNegative,
		AccountID: arg.FromAccountID,
	})
	if err != nil {
		return fmt.Errorf("insufficient funds or balance update failed: %w", err)
	}

	// Posting and accrual of the Recipient's balance
	_, err = qtx.CreatePosting(ctx, CreatePostingParams{
		TransactionID: txID.ID,
		AccountID:     pgtype.Int8{Int64: arg.ToAccountID, Valid: true},
		Amount:        amountPositive,
		CurrencyCode:  pgtype.Text{String: arg.CurrencyCode, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create posting for receiver: %w", err)
	}

	_, err = qtx.AddAccountBalance(ctx, AddAccountBalanceParams{
		Balance:   amountPositive,
		AccountID: arg.ToAccountID,
	})
	if err != nil {
		return fmt.Errorf("failed to add balance to receiver: %w", err)
	}
	return tx.Commit(ctx)
}
