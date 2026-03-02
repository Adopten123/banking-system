package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDBPool *pgxpool.Pool

const migrationsPath = "file://../migrations"

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Starting PostgreSQL in Docker
	log.Println("Starting PostgreSQL container...")
	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("accounts_test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(10*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}

	defer func() {
		log.Println("Terminating PostgreSQL container...")
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate container: %v", err)
		}
	}()

	// Getting connection str for DB
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	// 3. Make migrations
	migrator, err := migrate.New(migrationsPath, connStr)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Migrations applied successfully!")

	// 4. Init connection pools
	testDBPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to create db pool: %v", err)
	}
	defer testDBPool.Close()

	// 5. Running tests
	log.Println("Running tests...")
	exitCode := m.Run()

	// 6. Exit with the test results
	os.Exit(exitCode)
}
