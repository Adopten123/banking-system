package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	_ "github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepo struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewAccountRepo(db *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{
		db:      db,
		queries: New(db),
	}
}

func (r *AccountRepo) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}

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

func (r *AccountRepo) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Account, error) {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(publicID.String()); err != nil {
		return nil, fmt.Errorf("invalid public_id: %w", err)
	}

	dbAcc, err := r.queries.GetAccountByPublicID(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("database execution failed: %w", err)
	}

	return &domain.Account{
		ID:           dbAcc.ID,
		PublicID:     publicID,
		UserID:       dbAcc.UserID.Bytes,
		TypeID:       dbAcc.TypeID.Int32,
		StatusID:     dbAcc.StatusID.Int32,
		CurrencyCode: dbAcc.CurrencyCode.String,
		Name:         dbAcc.Name.String,
		Version:      dbAcc.Version.Int32,
		CreatedAt:    dbAcc.CreatedAt.Time,
		UpdatedAt:    dbAcc.UpdatedAt.Time,
	}, nil
}
