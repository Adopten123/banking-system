package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *socialService) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	return s.repo.Follow(ctx, followerID, followingID)
}

func (s *socialService) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	return s.repo.Unfollow(ctx, followerID, followingID)
}