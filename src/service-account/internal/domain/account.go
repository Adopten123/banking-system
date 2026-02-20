package domain

import "context"

type Account struct {
	ID      int
	Balance float64
}

type AccountRepository interface {
	Ping(ctx context.Context) error
}

type AccountService interface {
	CheckHealth(ctx context.Context) string
}
