package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *socialService) ToggleLike(ctx context.Context, postID int64, userID uuid.UUID, isLike bool) error {
	if isLike {
		return s.repo.Like(ctx, postID, userID)
	}
	return s.repo.Unlike(ctx, postID, userID)
}