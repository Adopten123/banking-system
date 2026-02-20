package postgres

import (
	"context"

	_ "github.com/Adopten123/banking-system/service-account/internal/domain"
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
