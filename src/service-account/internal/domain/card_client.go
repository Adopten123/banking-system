package domain

import "context"

// CardVaultClient — service-safe interface
type CardVaultClient interface {
	IssueCard(ctx context.Context, params IssueCardParams) (IssuedCardData, error)
	DeleteCardData(ctx context.Context, tokenID string) error
	VerifyCard(ctx context.Context, input VerifyCardInput) (bool, string, error)

	GetCardDetails(ctx context.Context, tokenID string) (*CardDetails, error)
	UpdateCardStatus(ctx context.Context, tokenID string, status string) error

	SetPin(ctx context.Context, tokenID string, pin string) error
	VerifyPin(ctx context.Context, tokenID string, pin string) (bool, error)
	GetTokenByPan(ctx context.Context, pan string) (string, error)
}
