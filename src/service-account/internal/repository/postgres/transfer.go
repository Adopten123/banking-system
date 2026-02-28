package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *AccountRepo) TransferTx(ctx context.Context, params domain.TransferParams) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
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
		return fmt.Errorf("failed to lock account 1: %w", err)
	}

	acc2, err := qtx.GetAccountForUpdate(ctx, account2ID)
	if err != nil {
		return fmt.Errorf("failed to lock account 2: %w", err)
	}

	// Business Checks. Verifying Statuses and Currencies
	var senderStatus, receiverStatus int32
	var senderCurrency, receiverCurrency string
	var senderBalance float64

	if params.FromAccountID == account1ID {
		senderStatus = acc1.StatusID.Int32
		receiverStatus = acc2.StatusID.Int32

		senderCurrency = acc1.CurrencyCode.String
		receiverCurrency = acc2.CurrencyCode.String

		bal, _ := acc1.Balance.Float64Value()
		senderBalance = bal.Float64
	} else {
		senderStatus = acc2.StatusID.Int32
		receiverStatus = acc1.StatusID.Int32

		senderCurrency = acc2.CurrencyCode.String
		receiverCurrency = acc1.CurrencyCode.String

		bal, _ := acc2.Balance.Float64Value()
		senderBalance = bal.Float64
	}

	if senderStatus != 1 {
		return domain.ErrAccountInactive
	}
	if receiverStatus != 1 {
		return domain.ErrAccountInactive
	}

	if senderCurrency != receiverCurrency {
		return errors.New("cross-currency transfers are not supported yet")
	}

	transferAmount, _ := strconv.ParseFloat(params.AmountStr, 64)

	if senderBalance < transferAmount {
		return domain.ErrInsufficientFunds
	}
	// Prepare Sum
	var amountPositive, amountNegative pgtype.Numeric

	if err := amountPositive.Scan(params.AmountStr); err != nil {
		return fmt.Errorf("failed to parse positive amount: %w", err)
	}

	if err := amountNegative.Scan("-" + params.AmountStr); err != nil {
		return fmt.Errorf("failed to parse negative amount: %w", err)
	}

	// ДЕБАГ: Выведем в консоль то, что реально полетит в базу
	fmt.Printf("DEBUG TRANSFER: SenderID=%d, ReceiverID=%d, Pos=%v, Neg=%v\n",
		params.FromAccountID, params.ToAccountID, amountPositive.Int, amountNegative.Int)

	// Making transaction
	//(category_id = 3 - transfer, status_id = 2 - posted)
	txID, err := qtx.CreateTransaction(ctx, CreateTransactionParams{
		ID:             pgtype.UUID{Bytes: uuid.New(), Valid: true},
		IdempotencyKey: pgtype.Text{String: params.IdempotencyKey, Valid: true},
		CategoryID:     pgtype.Int4{Int32: 3, Valid: true},
		StatusID:       pgtype.Int4{Int32: 2, Valid: true},
		Description:    pgtype.Text{String: params.Description, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Posting and writing off the Sender's balance
	_, err = qtx.CreatePosting(ctx, CreatePostingParams{
		TransactionID: txID.ID,
		AccountID:     pgtype.Int8{Int64: params.FromAccountID, Valid: true},
		Amount:        amountNegative,
		CurrencyCode:  pgtype.Text{String: params.CurrencyCode, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create posting for sender: %w", err)
	}

	_, err = qtx.AddAccountBalance(ctx, AddAccountBalanceParams{
		Balance:   amountNegative,
		AccountID: params.FromAccountID,
	})
	if err != nil {
		return fmt.Errorf("insufficient funds or balance update failed: %w", err)
	}

	// Posting and accrual of the Recipient's balance
	_, err = qtx.CreatePosting(ctx, CreatePostingParams{
		TransactionID: txID.ID,
		AccountID:     pgtype.Int8{Int64: params.ToAccountID, Valid: true},
		Amount:        amountPositive,
		CurrencyCode:  pgtype.Text{String: params.CurrencyCode, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create posting for receiver: %w", err)
	}

	_, err = qtx.AddAccountBalance(ctx, AddAccountBalanceParams{
		Balance:   amountPositive,
		AccountID: params.ToAccountID,
	})

	if err != nil {
		return fmt.Errorf("failed to add balance to receiver: %w", err)
	}
	return tx.Commit(ctx)
}
