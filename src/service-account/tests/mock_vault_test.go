package tests

import (
	"context"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
)

type MockCardVaultClient struct {
	IssueCardFunc        func(ctx context.Context, params domain.IssueCardParams) (domain.IssuedCardData, error)
	GetCardDetailsFunc   func(ctx context.Context, tokenID string) (*domain.CardDetails, error)
	UpdateCardStatusFunc func(ctx context.Context, tokenID string, status string) error
	SetPinFunc           func(ctx context.Context, tokenID string, pin string) error
	VerifyPinFunc        func(ctx context.Context, tokenID string, pin string) (bool, error)
	DeleteCardDataFunc   func(ctx context.Context, tokenID string) error
	VerifyCardFunc       func(ctx context.Context, input domain.VerifyCardInput) (bool, string, error)
	GetTokenByPanFunc    func(ctx context.Context, pan string) (string, error)
}

func (m *MockCardVaultClient) IssueCard(ctx context.Context, params domain.IssueCardParams) (domain.IssuedCardData, error) {
	if m.IssueCardFunc != nil {
		return m.IssueCardFunc(ctx, params)
	}
	return domain.IssuedCardData{}, nil
}

func (m *MockCardVaultClient) GetCardDetails(ctx context.Context, tokenID string) (*domain.CardDetails, error) {
	if m.GetCardDetailsFunc != nil {
		return m.GetCardDetailsFunc(ctx, tokenID)
	}
	return nil, nil
}

func (m *MockCardVaultClient) UpdateCardStatus(ctx context.Context, tokenID string, status string) error {
	if m.UpdateCardStatusFunc != nil {
		return m.UpdateCardStatusFunc(ctx, tokenID, status)
	}
	return nil
}

func (m *MockCardVaultClient) SetPin(ctx context.Context, tokenID string, pin string) error {
	if m.SetPinFunc != nil {
		return m.SetPinFunc(ctx, tokenID, pin)
	}
	return nil
}

func (m *MockCardVaultClient) VerifyPin(ctx context.Context, tokenID string, pin string) (bool, error) {
	if m.VerifyPinFunc != nil {
		return m.VerifyPinFunc(ctx, tokenID, pin)
	}
	return true, nil
}

func (m *MockCardVaultClient) DeleteCardData(ctx context.Context, tokenID string) error {
	if m.DeleteCardDataFunc != nil {
		return m.DeleteCardDataFunc(ctx, tokenID)
	}
	return nil
}

func (m *MockCardVaultClient) VerifyCard(ctx context.Context, input domain.VerifyCardInput) (bool, string, error) {
	if m.VerifyCardFunc != nil {
		return m.VerifyCardFunc(ctx, input)
	}
	return true, "", nil
}

func (m *MockCardVaultClient) GetTokenByPan(ctx context.Context, pan string) (string, error) {
	if m.GetTokenByPanFunc != nil {
		return m.GetTokenByPanFunc(ctx, pan)
	}
	return "", nil
}