package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *socialRepository) Like(ctx context.Context, postID int64, userID uuid.UUID) error {
	err := r.q.AddPostLike(ctx, AddPostLikeParams{
		PostID: postID,
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err == nil {
		_ = r.q.IncrementPostLikes(ctx, postID)
	}
	return err
}

func (r *socialRepository) Unlike(ctx context.Context, postID int64, userID uuid.UUID) error {
	err := r.q.RemovePostLike(ctx, RemovePostLikeParams{
		PostID: postID,
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err == nil {
		_ = r.q.DecrementPostLikes(ctx, postID)
	}
	return err
}
