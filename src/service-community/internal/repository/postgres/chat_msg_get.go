package postgres

import (
	"context"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *chatRepository) GetChatMessages(
	ctx context.Context,
	chatID uuid.UUID,
	limit, offset int32,
) ([]domain.Message, error) {

	params := GetChatMessagesParams{
		ChatID: pgtype.UUID{Bytes: chatID, Valid: true},
		Limit:  limit,
		Offset: offset,
	}

	rows, err := r.q.GetChatMessages(ctx, params)
	if err != nil {
		return nil, err
	}

	var messages []domain.Message
	for _, row := range rows {
		messages = append(messages, *mapMessageToDomain(row))
	}
	return messages, nil
}
