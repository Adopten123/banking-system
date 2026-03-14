package broker

import (
	"context"
	"github.com/Adopten123/banking-system/service-account/internal/domain"
)

func (p *RabbitMQPublisher) PublishTransferCreated(ctx context.Context, event domain.TransferCreatedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishAccountCreated(ctx context.Context, event domain.AccountCreatedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishAccountStatusChanged(ctx context.Context, event domain.AccountStatusChangedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishDepositCompleted(ctx context.Context, event domain.DepositCompletedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishWithdrawalCompleted(ctx context.Context, event domain.WithdrawalCompletedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishCreditLimitChanged(ctx context.Context, event domain.CreditLimitChangedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishCardIssued(ctx context.Context, event domain.CardIssuedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishCardStatusChanged(ctx context.Context, event domain.CardStatusChangedEvent) error {
	return p.publishMessage(ctx, event)
}
