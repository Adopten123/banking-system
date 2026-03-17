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

func (r *chatRepository) GetChatMemberIDs(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.q.GetChatMemberIDs(ctx, pgtype.UUID{Bytes: chatID, Valid: true})
	if err != nil {
		return nil, err
	}

	var members []uuid.UUID
	for _, row := range rows {
		members = append(members, row.Bytes)
	}
	return members, nil
}