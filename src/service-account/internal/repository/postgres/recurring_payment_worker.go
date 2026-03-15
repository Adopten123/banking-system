package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (r *AccountRepo) GetDueRecurringPayments(ctx context.Context, limit int32) ([]domain.RecurringPayment, error) {
	pgPayments, err := r.queries.GetDueRecurringPayments(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get due recurring payments: %w", err)
	}

	domainPayments := make([]domain.RecurringPayment, len(pgPayments))

	for i, p := range pgPayments {
		var id, srcID, destID uuid.UUID
		_ = id.UnmarshalBinary(p.ID.Bytes[:])
		_ = srcID.UnmarshalBinary(p.SourceID.Bytes[:])

		if p.DestinationID.Valid {
			_ = destID.UnmarshalBinary(p.DestinationID.Bytes[:])
		}

		var amount decimal.Decimal
		val, _ := p.Amount.Value()
		_ = amount.Scan(val)

		domainPayments[i] = domain.RecurringPayment{
			ID:                id,
			SourceTypeID:      p.SourceTypeID,
			SourceID:          srcID,
			DestinationTypeID: p.DestinationTypeID.Int32,
			DestinationID:     destID,
			Amount:            amount,
			CurrencyCode:      p.CurrencyCode.String,
			CategoryID:        p.CategoryID.Int32,
			CronExpression:    p.CronExpression,
			NextExecutionTime: p.NextExecutionTime.Time,
			IsActive:          p.IsActive.Bool,
			Description:       p.Description.String,
		}
	}
	return domainPayments, nil
}

func (r *AccountRepo) UpdateRecurringPaymentNextRun(
	ctx context.Context,
	id uuid.UUID,
	nextRun time.Time,
) error {
	var pgID pgtype.UUID
	_ = pgID.Scan(id.String())

	params := UpdateRecurringPaymentNextRunParams{
		ID:                pgID,
		NextExecutionTime: pgtype.Timestamptz{Time: nextRun, Valid: true},
	}

	err := r.queries.UpdateRecurringPaymentNextRun(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update next run time: %w", err)
	}
	return nil
}
