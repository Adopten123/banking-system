package postgres

import (
	"context"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *postRepository) Create(ctx context.Context, input domain.CreatePostInput) (*domain.Post, error) {
	params := CreatePostParams{
		AuthorID:         pgtype.UUID{Bytes: input.AuthorID, Valid: true},
		TypeID:           input.TypeID,
		Content:          input.Content,
		MediaAttachments: input.MediaAttachments,
		Status:           input.Status,
	}

	if input.RelatedAssetTicker != nil {
		params.RelatedAssetTicker = pgtype.Text{String: *input.RelatedAssetTicker, Valid: true}
	}

	row, err := r.q.CreatePost(ctx, params)
	if err != nil {
		return nil, err
	}

	var ticker *string
	if row.RelatedAssetTicker.Valid {
		ticker = &row.RelatedAssetTicker.String
	}

	return &domain.Post{
		ID:                 row.ID,
		AuthorID:           row.AuthorID.Bytes,
		TypeID:             row.TypeID,
		Content:            row.Content,
		MediaAttachments:   row.MediaAttachments,
		RelatedAssetTicker: ticker,
		Status:             row.Status,
		LikesCount:         row.LikesCount,
		CommentsCount:      row.CommentsCount,
		IsPinned:           row.IsPinned,
		IsEdited:           row.IsEdited,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
	}, nil
}