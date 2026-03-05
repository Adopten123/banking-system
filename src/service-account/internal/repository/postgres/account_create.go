package postgres

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	_ "github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *AccountRepo) Create(ctx context.Context, acc *domain.Account) (*domain.Account, error) {
	// Open transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	dbAcc, err := qtx.CreateAccount(ctx, CreateAccountParams{
		PublicID:     pgtype.UUID{Bytes: acc.PublicID, Valid: true},
		UserID:       pgtype.UUID{Bytes: acc.UserID, Valid: true},
		TypeID:       pgtype.Int4{Int32: acc.TypeID, Valid: true},
		StatusID:     pgtype.Int4{Int32: acc.StatusID, Valid: true},
		CurrencyCode: pgtype.Text{String: acc.CurrencyCode, Valid: true},
		Name:         pgtype.Text{String: acc.Name, Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to insert account: %w", err)
	}

	err = qtx.CreateAccountBalance(ctx, dbAcc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert account balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	acc.ID = dbAcc.ID
	acc.Version = dbAcc.Version.Int32
	acc.CreatedAt = dbAcc.CreatedAt.Time
	acc.UpdatedAt = dbAcc.UpdatedAt.Time

	return acc, nil
}
