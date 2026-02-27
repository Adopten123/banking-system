package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *AccountRepo) GetTransactions(
	ctx context.Context,
	accountID int64,
	limit, offset int32,
	startDate, endDate *time.Time,
) ([]domain.TransactionHistory, error) {

	var startPg, endPg pgtype.Timestamp
	if startDate != nil {
		startPg = pgtype.Timestamp{Time: *startDate, Valid: true}
	}
	if endDate != nil {
		endPg = pgtype.Timestamp{Time: *endDate, Valid: true}
	}

	params := GetAccountTransactionsParams{
		AccountID: pgtype.Int8{Int64: accountID, Valid: true},
		StartDate: startPg,
		EndDate:   endPg,
		Limit:     limit,
		Offset:    offset,
	}

	rows, err := r.queries.GetAccountTransactions(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	var history []domain.TransactionHistory
	for _, row := range rows {
		history = append(history, domain.TransactionHistory{
			TransactionID: row.TransactionID.Bytes,
			CategoryID:    row.CategoryID.Int32,
			StatusID:      row.StatusID.Int32,
			Description:   row.Description.String,
			Amount:        row.AmountStr,
			CurrencyCode:  row.CurrencyCode.String,
			CreatedAt:     row.CreatedAt.Time,
		})
	}

	return history, nil
}
