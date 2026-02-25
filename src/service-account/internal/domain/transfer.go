package domain

type TransferParams struct {
	FromAccountID  int64
	ToAccountID    int64
	AmountStr      string
	CurrencyCode   string
	IdempotencyKey string
	Description    string
}
