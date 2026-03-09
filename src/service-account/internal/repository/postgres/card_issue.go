package postgres

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *AccountRepo) CreateCard(ctx context.Context, card *domain.Card) error {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(card.ID.String()); err != nil {
		return fmt.Errorf("invalid card id: %w", err)
	}

	params := CreateCardParams{
		ID:         pgUUID,
		AccountID:  pgtype.Int8{Int64: card.AccountID, Valid: true},
		PanMask:    pgtype.Text{String: card.PANMask, Valid: true},
		ExpiryDate: pgtype.Date{Time: card.Expiry, Valid: true},
		IsVirtual:  pgtype.Bool{Bool: card.IsVirtual, Valid: true},
		Status:     pgtype.Text{String: card.Status, Valid: true},
	}

	createdCard, err := r.queries.CreateCard(ctx, params)
	if err != nil {
		return fmt.Errorf("database execution failed: %w", err)
	}

	card.CreatedAt = createdCard.CreatedAt.Time

	return nil
}
