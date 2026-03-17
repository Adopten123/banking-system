package postgres

import (
	"context"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *chatRepository) CreateChat(ctx context.Context, typeID int32, title, avatarURL *string) (*domain.Chat, error) {
	params := CreateChatParams{
		TypeID: typeID,
	}
	if title != nil {
		params.Title = pgtype.Text{String: *title, Valid: true}
	}
	if avatarURL != nil {
		params.AvatarUrl = pgtype.Text{String: *avatarURL, Valid: true}
	}

	row, err := r.q.CreateChat(ctx, params)
	if err != nil {
		return nil, err
	}

	chat := &domain.Chat{
		ID:            row.ID.Bytes,
		TypeID:        row.TypeID,
		LastMessageAt: row.LastMessageAt.Time,
		CreatedAt:     row.CreatedAt.Time,
	}
	if row.Title.Valid {
		chat.Title = &row.Title.String
	}
	if row.AvatarUrl.Valid {
		chat.AvatarURL = &row.AvatarUrl.String
	}
	return chat, nil
}