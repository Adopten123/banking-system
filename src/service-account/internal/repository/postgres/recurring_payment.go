package postgres

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *AccountRepo) CreateRecurringPayment(ctx context.Context, p domain.RecurringPayment) error {
	var pgID pgtype.UUID
	_ = pgID.Scan(p.ID.String())

	var pgSourceID pgtype.UUID
	_ = pgSourceID.Scan(p.SourceID.String())

	var pgDestTypeID pgtype.Int4
	var pgDestID pgtype.UUID
	if p.DestinationTypeID != 0 && p.DestinationID != uuid.Nil {
		pgDestTypeID = pgtype.Int4{Int32: p.DestinationTypeID, Valid: true}
		_ = pgDestID.Scan(p.DestinationID.String())
	}

	var pgAmount pgtype.Numeric
	_ = pgAmount.Scan(p.Amount.String())

	var pgCategory pgtype.Int4
	if p.CategoryID != 0 {
		pgCategory = pgtype.Int4{Int32: p.CategoryID, Valid: true}
	}

	var pgDesc pgtype.Text
	if p.Description != "" {
		pgDesc = pgtype.Text{String: p.Description, Valid: true}
	}

	var pgNextTime pgtype.Timestamptz
	_ = pgNextTime.Scan(p.NextExecutionTime)

	_, err := r.queries.CreateRecurringPayment(ctx, CreateRecurringPaymentParams{
		ID:                pgID,
		SourceTypeID:      p.SourceTypeID,
		SourceID:          pgSourceID,
		DestinationTypeID: pgDestTypeID,
		DestinationID:     pgDestID,
		Amount:            pgAmount,
		CurrencyCode:      pgtype.Text{String: p.CurrencyCode, Valid: true},
		CategoryID:        pgCategory,
		CronExpression:    p.CronExpression,
		NextExecutionTime: pgNextTime,
		Description:       pgDesc,
	})

	if err != nil {
		return fmt.Errorf("failed to create recurring payment: %w", err)
	}

	return nil
}

func (r *AccountRepo) CancelRecurringPayment(ctx context.Context, id uuid.UUID) error {
	var pgID pgtype.UUID
	_ = pgID.Scan(id.String())

	err := r.queries.UpdateRecurringPaymentStatus(ctx, UpdateRecurringPaymentStatusParams{
		ID:       pgID,
		IsActive: pgtype.Bool{Bool: false, Valid: true},
	})

	if err != nil {
		return fmt.Errorf("failed to cancel recurring payment: %w", err)
	}
	return nil
}
