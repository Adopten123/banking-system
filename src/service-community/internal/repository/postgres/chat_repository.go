package postgres

import (
	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/google/uuid"
)

type chatRepository struct {
	q *Queries
}

func NewChatRepository(q *Queries) domain.ChatRepository {
	return &chatRepository{q: q}
}

func mapMessageToDomain(row Message) *domain.Message {
	msg := &domain.Message{
		ID:         row.ID,
		ChatID:     row.ChatID.Bytes,
		SenderID:   row.SenderID.Bytes,
		IsTransfer: row.IsTransfer,
		IsEdited:   row.IsEdited,
		CreatedAt:  row.CreatedAt.Time,
	}

	if row.Content.Valid {
		msg.Content = &row.Content.String
	}
	if row.ReplyToMessageID.Valid {
		msg.ReplyToMessageID = &row.ReplyToMessageID.Int64
	}
	if row.TransferAmount.Valid {
		val, _ := row.TransferAmount.Value()
		if strVal, ok := val.(string); ok {
			msg.TransferAmount = &strVal
		}
	}
	if row.TransferCurrency.Valid {
		msg.TransferCurrency = &row.TransferCurrency.String
	}
	if row.IdempotencyKey.Valid {
		msg.IdempotencyKey = &row.IdempotencyKey.String
	}
	if row.TransferStatus.Valid {
		msg.TransferStatus = &row.TransferStatus.String
	}
	if row.TransferTransactionID.Valid {
		uuidBytes := row.TransferTransactionID.Bytes
		msg.TransferTransactionID = (*uuid.UUID)(&uuidBytes)
	}
	if row.DeletedAt.Valid {
		msg.DeletedAt = &row.DeletedAt.Time
	}

	return msg
}