package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// createTestUser making unique user for test
func createTestUser(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	// Generate unique uuid
	userID := uuid.New()

	// Generate unique name
	username := fmt.Sprintf("test_user_%s", userID.String()[:8])

	// Making SQL query
	query := `
		INSERT INTO "users" ("id", "username") 
		VALUES ($1, $2) 
		ON CONFLICT ("id") DO NOTHING;
	`

	_, err := pool.Exec(context.Background(), query, userID, username)
	require.NoError(t, err, "failed to create test user in database")

	return userID
}

// createTestAccount - make test account with balance == 0
// return account_id
func createTestAccount(t *testing.T, pool *pgxpool.Pool, userID uuid.UUID, currencyCode string) uuid.UUID {
	var internalID int64
	publicID := uuid.New()

	queryAccount := `
		INSERT INTO accounts (public_id, user_id, type_id, status_id, currency_code, name) 
		VALUES ($1, $2, 1, 1, $3, 'Test Account') 
		RETURNING id;
	`

	err := pool.QueryRow(context.Background(), queryAccount, publicID, userID, currencyCode).Scan(&internalID)
	require.NoError(t, err, "failed to insert into accounts table")

	queryBalance := `
		INSERT INTO account_balances (account_id, balance, credit_limit) 
		VALUES ($1, 0, 0);
	`
	_, err = pool.Exec(context.Background(), queryBalance, internalID)
	require.NoError(t, err, "failed to insert into account_balances table")

	return publicID
}

// getAccountBalance - get balance of account
func getAccountBalance(t *testing.T, pool *pgxpool.Pool, publicID uuid.UUID) decimal.Decimal {
	var balance decimal.Decimal

	query := `
		SELECT ab.balance 
		FROM account_balances ab
		JOIN accounts a ON a.id = ab.account_id
		WHERE a.public_id = $1;
	`
	err := pool.QueryRow(context.Background(), query, publicID).Scan(&balance)
	require.NoError(t, err, "failed to get account balance")

	return balance
}

// setAccountBalance - set balance of account
func setAccountBalance(t *testing.T, pool *pgxpool.Pool, publicID uuid.UUID, amount decimal.Decimal) {
	query := `
		UPDATE account_balances ab
		SET balance = $1
		FROM accounts a
		WHERE a.id = ab.account_id AND a.public_id = $2;
	`
	_, err := pool.Exec(context.Background(), query, amount, publicID)
	require.NoError(t, err, "failed to set account balance")
}
