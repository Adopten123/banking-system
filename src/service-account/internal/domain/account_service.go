package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountService interface {
	CheckHealth(ctx context.Context) string

	CreateAccount(ctx context.Context, params CreateAccountInput) (*Account, error)

	Deposit(ctx context.Context, publicID uuid.UUID, input ServiceDepositInput) error
	Withdraw(ctx context.Context, publicID uuid.UUID, amount decimal.Decimal, idempotencyKey string) (*WithdrawResponse, error)
	Transfer(ctx context.Context, input TransferInput) error

	GetAccount(ctx context.Context, publicID uuid.UUID) (*Account, error)
	GetAccountBalance(ctx context.Context, publicID uuid.UUID) (decimal.Decimal, error)
	GetAccountTransactions(ctx context.Context, publicID uuid.UUID, params TransactionFilter) ([]TransactionHistory, error)

	UpdateCreditLimit(ctx context.Context, publicID uuid.UUID, limitStr string) error
	BlockAccount(ctx context.Context, publicID uuid.UUID) error
	CloseAccount(ctx context.Context, publicID uuid.UUID) error
}
