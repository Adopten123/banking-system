package tests

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

type MockExchangeClient struct{}

func (m *MockExchangeClient) GetRate(ctx context.Context, base, target string) (decimal.Decimal, error) {
	if target == "EUR" {
		return decimal.Zero, fmt.Errorf("exchange service is down")
	}

	if base == "RUB" && target == "USD" {
		return decimal.RequireFromString("0.0108"), nil
	}

	return decimal.NewFromInt(1), nil
}