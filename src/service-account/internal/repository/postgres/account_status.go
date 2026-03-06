package postgres

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (r *AccountRepo) UpdateStatus(ctx context.Context, accountID int64, statusID int32) error {
	err := r.queries.UpdateAccountStatus(ctx, UpdateAccountStatusParams{
		StatusID: pgtype.Int4{Int32: statusID, Valid: true},
		ID:       accountID,
	})

	if err != nil {
		return fmt.Errorf("failed to update account status: %w", err)
	}

	return nil
}

func (r *AccountRepo) CloseAccountTx(ctx context.Context, accountID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Blocked balance and get it
	balancePg, err := qtx.GetBalanceForUpdate(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to lock balance for update: %w", err)
	}

	var currentBalance decimal.Decimal
	if val, err := balancePg.Balance.Value(); err == nil && val != nil {
		_ = currentBalance.Scan(val)
	}

	if !currentBalance.IsZero() {
		return domain.ErrAccountHasBalance
	}

	// Change status to "Closed"
	//1, 'active',
	//2, 'frozen',
	//3, 'blocked',
	//4, 'closed'
	err = qtx.UpdateAccountStatus(ctx, UpdateAccountStatusParams{
		StatusID: pgtype.Int4{Int32: 4, Valid: true},
		ID:       accountID,
	})
	if err != nil {
		return fmt.Errorf("failed to set account status to closed: %w", err)
	}

	// Commit changes
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit close account transaction: %w", err)
	}

	return nil
}
