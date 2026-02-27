package domain

import (
	"context"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Ping(ctx context.Context) error
	Create(ctx context.Context, account *Account) (*Account, error)

	Deposit(ctx context.Context, params RepoDepositParams) error
	TransferTx(ctx context.Context, arg TransferParams) error

	GetByPublicID(ctx context.Context, publicID uuid.UUID) (*Account, error)
	GetTransactions(ctx context.Context, accountID int64, filter TransactionFilter) ([]TransactionHistory, error)
}

type AccountService interface {
	CheckHealth(ctx context.Context) string

	CreateAccount(ctx context.Context, params CreateAccountInput) (*Account, error)

	Deposit(ctx context.Context, publicID uuid.UUID, input ServiceDepositInput) error
	Transfer(ctx context.Context, input TransferInput) error

	GetAccount(ctx context.Context, publicID uuid.UUID) (*Account, error)
	GetAccountTransactions(ctx context.Context, publicID uuid.UUID, params TransactionFilter) ([]TransactionHistory, error)
}
