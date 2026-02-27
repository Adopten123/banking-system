package domain

// ServiceDepositInput - data for Deposit (Service Layer)
type ServiceDepositInput struct {
	AmountStr      string
	IdempotencyKey string
}

// RepoDepositParams - data for Deposit (Service Repository)
type RepoDepositParams struct {
	AccountID      int64
	AmountStr      string
	CurrencyCode   string
	IdempotencyKey string
}
