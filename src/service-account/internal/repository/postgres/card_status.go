package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *AccountRepo) UpdateCardStatus(ctx context.Context, cardID uuid.UUID, status string) error {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(cardID.String()); err != nil {
		return fmt.Errorf("invalid card id: %w", err)
	}

	params := UpdateCardStatusParams{
		ID: pgUUID,
		Status: pgtype.Text{
			String: status,
			Valid:  true,
		},
	}

	err := r.queries.UpdateCardStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("database execution failed: %w", err)
	}

	return nil
}
