package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *socialRepository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return r.q.FollowUser(ctx, FollowUserParams{
		FollowerID:  pgtype.UUID{Bytes: followerID, Valid: true},
		FollowingID: pgtype.UUID{Bytes: followingID, Valid: true},
	})
}

func (r *socialRepository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return r.q.UnfollowUser(ctx, UnfollowUserParams{
		FollowerID:  pgtype.UUID{Bytes: followerID, Valid: true},
		FollowingID: pgtype.UUID{Bytes: followingID, Valid: true},
	})
}
