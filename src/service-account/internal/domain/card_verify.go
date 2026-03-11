package domain

import "github.com/google/uuid"

// VerifyCardInput — input data for acquiring
type VerifyCardInput struct {
	PAN         string
	CVV         string
	ExpiryMonth int32
	ExpiryYear  int32
}

// VerifyCardResult — result of card checking
type VerifyCardResult struct {
	IsValid   bool      `json:"is_valid"`
	CardID    uuid.UUID `json:"card_id,omitempty"`
	AccountID int64     `json:"account_id,omitempty"`
}

type VerifyCardRequest struct {
	Pan         string `json:"pan"`
	Cvv         string `json:"cvv"`
	ExpiryMonth int32  `json:"expiry_month"`
	ExpiryYear  int32  `json:"expiry_year"`
}
