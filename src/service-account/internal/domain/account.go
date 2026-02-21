package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           int64
	PublicID     uuid.UUID
	UserID       uuid.UUID
	TypeID       int32
	StatusID     int32
	CurrencyCode string
	Name         string
	Version      int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AccountRepository interface {
	Ping(ctx context.Context) error
	Create(ctx context.Context, account *Account) (*Account, error)
}

type AccountService interface {
	CheckHealth(ctx context.Context) string
	CreateAccount(ctx context.Context, userID uuid.UUID, typeID int32, currencyCode, name string) (*Account, error)
}
