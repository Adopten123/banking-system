package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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