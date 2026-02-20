package postgres

import (
	_ "github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepo struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewAccountRepo() *AccountRepo {
	return &AccountRepo{}
}

func (r *AccountRepo) Ping() error {
	return nil
}
