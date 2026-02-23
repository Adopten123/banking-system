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
	GetByPublicID(ctx context.Context, publicID uuid.UUID) (*Account, error)
	Deposit(ctx context.Context, accountID int64, amountStr string, currencyCode string, idempotencyKey string) error
}

type AccountService interface {
	CheckHealth(ctx context.Context) string
	CreateAccount(ctx context.Context, userID uuid.UUID, typeID int32, currencyCode, name string) (*Account, error)
	GetAccount(ctx context.Context, publicID uuid.UUID) (*Account, error)
	Deposit(ctx context.Context, publicID uuid.UUID, amountStr string, idempotencyKey string) error
}
