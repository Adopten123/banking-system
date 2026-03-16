package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPostgres(ctx context.Context, connString string) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse DB config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return pool
}