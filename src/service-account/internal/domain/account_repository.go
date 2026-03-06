package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	Ping(ctx context.Context) error
	Create(ctx context.Context, account *Account) (*Account, error)

	DepositTx(ctx context.Context, params RepoDepositParams) (*DepositResult, error)
	TransferTx(ctx context.Context, arg TransferParams) (*TransferResult, error)
	WithdrawTx(ctx context.Context, publicID uuid.UUID, amount decimal.Decimal, idempotencyKey string) (*WithdrawResult, error)

	GetByPublicID(ctx context.Context, publicID uuid.UUID) (*Account, error)
	GetBalance(ctx context.Context, publicID uuid.UUID) (decimal.Decimal, error)
	GetTransactions(ctx context.Context, accountID int64, filter TransactionFilter) (*TransactionHistoryResult, error)

	GetCreditLimit(ctx context.Context, accountID int64) (string, error)
	UpdateCreditLimit(ctx context.Context, accountID int64, limitStr string) error

	UpdateStatus(ctx context.Context, accountID int64, statusID int32) error
	CloseAccountTx(ctx context.Context, accountID int64) error
}
