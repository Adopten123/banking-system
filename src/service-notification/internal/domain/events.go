package domain

import (
	"time"
)

type TransferCreatedEvent struct {
	TransactionID  string    `json:"transaction_id"`
	FromAccountID  int       `json:"from_account_id"`
	ToAccountID    int       `json:"to_account_id"`
	Amount         string    `json:"amount"`
	Currency       string    `json:"currency"`
	IdempotencyKey string    `json:"idempotency_key"`
	Timestamp      time.Time `json:"timestamp"`
}

type AccountCreatedEvent struct {
	AccountID int       `json:"account_id"`
	PublicID  string    `json:"public_id"`
	UserID    string    `json:"user_id"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

type AccountStatusChangedEvent struct {
	AccountID int       `json:"account_id"`
	OldStatus int       `json:"old_status"`
	NewStatus int       `json:"new_status"`
	Timestamp time.Time `json:"timestamp"`
}

type DepositCompletedEvent struct {
	TransactionID string    `json:"transaction_id"`
	AccountID     int       `json:"account_id"`
	Amount        string    `json:"amount"`
	Currency      string    `json:"currency"`
	Timestamp     time.Time `json:"timestamp"`
}

type WithdrawalCompletedEvent struct {
	TransactionID string    `json:"transaction_id"`
	AccountID     int       `json:"account_id"`
	Amount        string    `json:"amount"`
	Currency      string    `json:"currency"`
	Timestamp     time.Time `json:"timestamp"`
}
