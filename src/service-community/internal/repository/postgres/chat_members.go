package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *chatRepository) AddChatMember(ctx context.Context, chatID, userID uuid.UUID, role string) error {
	return r.q.AddChatMember(ctx, AddChatMemberParams{
		ChatID: pgtype.UUID{Bytes: chatID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Role:   role,
	})
}
