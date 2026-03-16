package service

import (
	"context"
	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"log"
)

func (s *postService) CreatePost(ctx context.Context, input domain.CreatePostInput) (*domain.Post, error) {
	if input.Content == "" {
		return nil, domain.ErrEmptyContent
	}

	if input.Status == "" {
		input.Status = "published"
	}

	// TODO: Add Investment/Exchange Service

	if input.RelatedAssetTicker != nil {
		log.Printf("Validating asset ticker: %s", *input.RelatedAssetTicker)
	}

	post, err := s.repo.Create(ctx, input)
	if err != nil {
		log.Printf("failed to create post in db: %v", err)
		return nil, err
	}

	// TODO: Add RabbitMQ

	return post, nil
}
