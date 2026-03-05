package postgres

import (
	"context"
	"errors"

	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (r *AccountRepo) GetBalance(ctx context.Context, publicID uuid.UUID) (decimal.Decimal, error) {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(publicID.String()); err != nil {
		return decimal.Zero, fmt.Errorf("invalid public_id: %w", err)
	}

	balanceStr, err := r.queries.GetAccountBalanceByPublicID(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return decimal.Zero, domain.ErrAccountNotFound
		}
		return decimal.Zero, fmt.Errorf("database execution failed: %w", err)
	}

	balance, err := decimal.NewFromString(balanceStr)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to parse balance from db: %w", err)
	}

	return balance, nil
}