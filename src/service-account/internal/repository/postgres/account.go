package postgres

import _ "github.com/Adopten123/banking-system/service-account/internal/domain"

type AccountRepo struct {
	// db *pgxpool.Pool
}

func NewAccountRepo() *AccountRepo {
	return &AccountRepo{}
}

func (r *AccountRepo) Ping() error {
	return nil
}
