package postgres

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *AccountRepo) GetTransactions(ctx context.Context, accountID int64) ([]domain.TransactionHistory, error) {
	pgAccID := pgtype.Int8{
		Int64: accountID,
		Valid: true,
	}

	rows, err := r.queries.GetAccountTransactions(ctx, pgAccID)
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
