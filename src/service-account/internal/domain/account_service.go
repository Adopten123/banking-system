package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountService interface {
	CheckHealth(ctx context.Context) string

	CreateAccount(ctx context.Context, params CreateAccountInput) (*Account, error)

	Deposit(ctx context.Context, input ServiceDepositInput) (*DepositResult, error)
	Withdraw(ctx context.Context, input ServiceWithdrawInput) (*WithdrawResult, error)

	Transfer(ctx context.Context, input TransferInput) (*TransferResult, error)

	GetAccount(ctx context.Context, publicID uuid.UUID) (*Account, error)
	GetAccountBalance(ctx context.Context, publicID uuid.UUID) (decimal.Decimal, error)
	GetAccountTransactions(ctx context.Context, publicID uuid.UUID, params TransactionFilter) (*TransactionHistoryResult, error)

	UpdateCreditLimit(ctx context.Context, publicID uuid.UUID, limitStr string) error
	ActivateAccount(ctx context.Context, publicID uuid.UUID) error
	FreezeAccount(ctx context.Context, publicID uuid.UUID) error
	BlockAccount(ctx context.Context, publicID uuid.UUID) error
	CloseAccount(ctx context.Context, publicID uuid.UUID) error

	// ---- CARDS METHODS ----
	IssueCard(ctx context.Context, input IssueCardInput) (*Card, error)
	DeleteCard(ctx context.Context, cardID uuid.UUID) error
	VerifyCardForPayment(ctx context.Context, input VerifyCardInput) (*VerifyCardResult, error)
	GetCardDetails(ctx context.Context, cardID uuid.UUID) (*CardDetails, error)
	GetAccountCards(ctx context.Context, accountPublicID uuid.UUID) ([]*Card, error)

	UpdateCardStatus(ctx context.Context, cardID uuid.UUID, status string) error

	SetCardPin(ctx context.Context, cardID uuid.UUID, pin string) error
	VerifyCardPin(ctx context.Context, cardID uuid.UUID, pin string) (bool, error)
}
