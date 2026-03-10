package domain

import "context"

// CardVaultClient — service-safe interface
type CardVaultClient interface {
	IssueCard(ctx context.Context, params IssueCardParams) (IssuedCardData, error)
	GetCardDetails(ctx context.Context, tokenID string) (*CardDetails, error)
}
