package domain

import (
	"context"
)

// IssueCardParams — params for card
type IssueCardParams struct {
	PaymentSystem string // "VISA", "MASTERCARD", "MIR"
	IsVirtual     bool
}

// IssuedCardData — safe sata, which service-safe returned
type IssuedCardData struct {
	TokenID     string
	PANMask     string
	ExpiryMonth int32
	ExpiryYear  int32
}

// CardVaultClient — service-safe interface
type CardVaultClient interface {
	IssueCard(ctx context.Context, params IssueCardParams) (IssuedCardData, error)
}
