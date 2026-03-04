package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *AccountRepo) UpdateCreditLimit(ctx context.Context, accountID int64, limitStr string) error {
	var limit pgtype.Numeric
	if err := limit.Scan(limitStr); err != nil {
		return fmt.Errorf("invalid credit limit format: %w", err)
	}

	err := r.queries.UpdateCreditLimit(ctx, UpdateCreditLimitParams{
		CreditLimit: limit,
		AccountID:   accountID,
	})
	if err != nil {
		return fmt.Errorf("failed to update credit limit: %w", err)
	}

	return nil
}
