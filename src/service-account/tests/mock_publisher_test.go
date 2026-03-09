package tests

import (
	"context"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
)

// MockPublisher implements domain.EventPublisher, save objects in memory for asserts
type MockPublisher struct {
	Events []domain.DomainEvent
}

func NewMockPublisher() *MockPublisher {
	return &MockPublisher{
		Events: make([]domain.DomainEvent, 0),
	}
}

func (m *MockPublisher) PublishTransferCreated(ctx context.Context, event domain.TransferCreatedEvent) error {
	m.Events = append(m.Events, event)
	return nil
}

func (m *MockPublisher) PublishAccountCreated(ctx context.Context, event domain.AccountCreatedEvent) error {
	m.Events = append(m.Events, event)
	return nil
}

func (m *MockPublisher) PublishAccountStatusChanged(ctx context.Context, event domain.AccountStatusChangedEvent) error {
	m.Events = append(m.Events, event)
	return nil
}

func (m *MockPublisher) PublishDepositCompleted(ctx context.Context, event domain.DepositCompletedEvent) error {
	m.Events = append(m.Events, event)
	return nil
}

func (m *MockPublisher) PublishWithdrawalCompleted(ctx context.Context, event domain.WithdrawalCompletedEvent) error {
	m.Events = append(m.Events, event)
	return nil
}

func (m *MockPublisher) PublishCreditLimitChanged(ctx context.Context, event domain.CreditLimitChangedEvent) error {
	m.Events = append(m.Events, event)
	return nil
}

// Clear - cleanses events before new tests
func (m *MockPublisher) Clear() {
	m.Events = make([]domain.DomainEvent, 0)
}
