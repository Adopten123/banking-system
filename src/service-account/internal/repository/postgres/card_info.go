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

// GetCardByID - find card in account_db by token
func (r *AccountRepo) GetCardByID(ctx context.Context, cardID uuid.UUID) (*domain.Card, error) {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(cardID.String()); err != nil {
		return nil, fmt.Errorf("invalid card id: %w", err)
	}

	row, err := r.queries.GetCardByID(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCardNotFound
		}
		return nil, fmt.Errorf("database execution failed: %w", err)
	}

	return &domain.Card{
		ID:        cardID,
		AccountID: row.AccountID.Int64,
		PANMask:   row.PanMask.String,
		Expiry:    row.ExpiryDate.Time,
		IsVirtual: row.IsVirtual.Bool,
		Status:    row.Status.String,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *AccountRepo) GetCardsByAccountID(ctx context.Context, accountID int64) ([]*domain.Card, error) {
	pgAccID := pgtype.Int8{
		Int64: accountID,
		Valid: true,
	}

	rows, err := r.queries.GetCardsByAccountID(ctx, pgAccID)
	if err != nil {
		return nil, fmt.Errorf("database execution failed: %w", err)
	}

	cards := make([]*domain.Card, 0, len(rows))

	for _, row := range rows {
		if !row.ID.Valid {
			continue
		}

		cardUUID := uuid.UUID(row.ID.Bytes)

		cards = append(cards, &domain.Card{
			ID:        cardUUID,
			AccountID: row.AccountID.Int64,
			PANMask:   row.PanMask.String,
			Expiry:    row.ExpiryDate.Time,
			IsVirtual: row.IsVirtual.Bool,
			Status:    row.Status.String,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return cards, nil
}
