package domain

import (
	"context"

	"github.com/google/uuid"
)

type SocialRepository interface {
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	Like(ctx context.Context, postID int64, userID uuid.UUID) error
	Unlike(ctx context.Context, postID int64, userID uuid.UUID) error
	CreateReport(ctx context.Context, input CreateReportInput) (*Report, error)
}

type SocialService interface {
	FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error
	UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error
	ToggleLike(ctx context.Context, postID int64, userID uuid.UUID, isLike bool) error
	FileReport(ctx context.Context, input CreateReportInput) (*Report, error)
}
