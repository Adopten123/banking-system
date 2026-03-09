package domain

import (
	"context"
	"github.com/shopspring/decimal"
)

type ExchangeRateClient interface {
	GetRate(ctx context.Context, baseCurrency, targetCurrency string) (decimal.Decimal, error)
}