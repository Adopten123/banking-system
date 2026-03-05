package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/Adopten123/banking-system/service-account/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepo struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewAccountRepo(db *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{
		db:      db,
		queries: New(db),
	}
}

func (r *AccountRepo) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}

func NewDBPool(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}
	log.Println("Successfully connected to PostgreSQL database")
	return pool, nil
}
