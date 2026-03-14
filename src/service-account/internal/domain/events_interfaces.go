package domain

import (
	"context"
)

type DomainEvent interface {
	EventName() string
}

type EventPublisher interface {
	PublishAccountCreated(ctx context.Context, event AccountCreatedEvent) error
	PublishAccountStatusChanged(ctx context.Context, event AccountStatusChangedEvent) error
	PublishCreditLimitChanged(ctx context.Context, event CreditLimitChangedEvent) error

	PublishTransferCreated(ctx context.Context, event TransferCreatedEvent) error
	PublishDepositCompleted(ctx context.Context, event DepositCompletedEvent) error
	PublishWithdrawalCompleted(ctx context.Context, event WithdrawalCompletedEvent) error

	PublishCardIssued(ctx context.Context, event CardIssuedEvent) error
	PublishCardStatusChanged(ctx context.Context, event CardStatusChangedEvent) error
}
