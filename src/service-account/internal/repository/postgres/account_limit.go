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

	err := r.queries.UpdateAccountCreditLimit(ctx, UpdateAccountCreditLimitParams{
		NewLimit:  limit,
		AccountID: accountID,
	})
	if err != nil {
		return fmt.Errorf("failed to update credit limit: %w", err)
	}

	return nil
}

func (r *AccountRepo) GetCreditLimit(ctx context.Context, accountID int64) (string, error) {
	limitStr, err := r.queries.GetCreditLimit(ctx, accountID)
	if err != nil {
		return "", fmt.Errorf("failed to get credit limit from db: %w", err)
	}

	return limitStr, nil
}
