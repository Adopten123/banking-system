package domain

import "github.com/google/uuid"

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

// IssueCardInput - data returned by HTTP-handler
type IssueCardInput struct {
	AccountPublicID uuid.UUID
	PaymentSystem   string
	IsVirtual       bool
}

type IssueCardRequest struct {
	PaymentSystem string `json:"payment_system"`
	IsVirtual     bool   `json:"is_virtual"`
}

