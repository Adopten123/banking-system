package domain

import (
	"context"
)

type PostRepository interface {
	Create(ctx context.Context, input CreatePostInput) (*Post, error)
	GetByID(ctx context.Context, id int64) (*Post, error)
}

type PostService interface {
	CreatePost(ctx context.Context, input CreatePostInput) (*Post, error)
}