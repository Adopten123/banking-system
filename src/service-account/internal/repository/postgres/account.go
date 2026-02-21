package postgres

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	_ "github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
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
	var pubID, usrID pgtype.UUID
	if err := pubID.Scan(acc.PublicID.String()); err != nil {
		return nil, fmt.Errorf("invalid public_id: %w", err)
	}
	if err := usrID.Scan(acc.UserID.String()); err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	arg := CreateAccountParams{
		PublicID:     pubID,
		UserID:       usrID,
		TypeID:       pgtype.Int4{Int32: acc.TypeID, Valid: true},
		StatusID:     pgtype.Int4{Int32: acc.StatusID, Valid: true},
		CurrencyCode: pgtype.Text{String: acc.CurrencyCode, Valid: true},
		Name:         pgtype.Text{String: acc.Name, Valid: true},
	}

	dbAcc, err := r.queries.CreateAccount(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create account in db: %w", err)
	}

	return &domain.Account{
		ID:           dbAcc.ID,
		PublicID:     acc.PublicID,
		UserID:       acc.UserID,
		TypeID:       dbAcc.TypeID.Int32,
		StatusID:     dbAcc.StatusID.Int32,
		CurrencyCode: dbAcc.CurrencyCode.String,
		Name:         dbAcc.Name.String,
		Version:      dbAcc.Version.Int32,
		CreatedAt:    dbAcc.CreatedAt.Time,
		UpdatedAt:    dbAcc.UpdatedAt.Time,
	}, nil
}

func (r *AccountRepo) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Account, error) {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(publicID.String()); err != nil {
		return nil, fmt.Errorf("invalid public_id: %w", err)
	}

	dbAcc, err := r.queries.GetAccountByPublicID(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
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
