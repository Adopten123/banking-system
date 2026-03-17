package postgres

import (
	"context"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *socialRepository) CreateReport(ctx context.Context, input domain.CreateReportInput) (*domain.Report, error) {
	row, err := r.q.CreateReport(ctx, CreateReportParams{
		ReporterID: pgtype.UUID{Bytes: input.ReporterID, Valid: true},
		TargetType: input.TargetType,
		TargetID:   input.TargetID,
		Reason:     input.Reason,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Report{
		ID:         row.ID,
		ReporterID: row.ReporterID.Bytes,
		TargetType: row.TargetType,
		TargetID:   row.TargetID,
		Reason:     row.Reason,
		Status:     row.Status,
		CreatedAt:  row.CreatedAt.Time,
	}, nil
}