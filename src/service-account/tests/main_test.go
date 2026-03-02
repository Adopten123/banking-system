package tests

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// 3. Init connection pools
	testDBPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to create db pool: %v", err)
	}
	defer testDBPool.Close()

	log.Println("Applying database schema...")

	schemaPath := filepath.Join("..", "internal", "repository", "postgres", "schema", "000001_init.up.sql")
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Failed to read schema file at %s: %v", schemaPath, err)
	}

	if _, err := testDBPool.Exec(ctx, string(schemaBytes)); err != nil {
		log.Fatalf("Failed to execute schema: %v", err)
	}
	log.Println("Schema applied successfully!")

	// 3. Make migrations
	log.Println("Applying data migrations...")
	migrator, err := migrate.New(migrationsPath, connStr)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Migrations applied successfully!")

	// 5. Running tests
	log.Println("Running tests...")
	exitCode := m.Run()

	// 6. Exit with the test results
	os.Exit(exitCode)
}
