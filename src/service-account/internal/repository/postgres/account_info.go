package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

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
