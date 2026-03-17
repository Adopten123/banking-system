package postgres

import (
	"context"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *chatRepository) CreateMessage(ctx context.Context, input domain.CreateMessageInput) (*domain.Message, error) {
	params := CreateMessageParams{
		ChatID:           pgtype.UUID{Bytes: input.ChatID, Valid: true},
		SenderID:         pgtype.UUID{Bytes: input.SenderID, Valid: true},
		IsTransfer:       input.IsTransfer,
	}

	if input.ReplyToMessageID != nil {
		params.ReplyToMessageID = pgtype.Int8{Int64: *input.ReplyToMessageID, Valid: true}
	}
	if input.Content != nil {
		params.Content = pgtype.Text{String: *input.Content, Valid: true}
	}
	if input.TransferCurrency != nil {
		params.TransferCurrency = pgtype.Text{String: *input.TransferCurrency, Valid: true}
	}
	if input.IdempotencyKey != nil {
		params.IdempotencyKey = pgtype.Text{String: *input.IdempotencyKey, Valid: true}
	}

	if input.TransferAmount != nil {
		params.TransferAmount.Scan(*input.TransferAmount)
	}

	row, err := r.q.CreateMessage(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapMessageToDomain(row), nil
}