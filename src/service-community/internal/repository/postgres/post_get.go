package postgres

import (
	"context"
	"github.com/Adopten123/banking-system/service-community/internal/domain"
)

func (r *postRepository) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	row, err := r.q.GetPostByID(ctx, id)
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