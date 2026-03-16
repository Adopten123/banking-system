package service

import "github.com/Adopten123/banking-system/service-community/internal/domain"

type postService struct {
	repo domain.PostRepository
}

// NewPostService инициализирует сервис бизнес-логики
func NewPostService(repo domain.PostRepository) domain.PostService {
	return &postService{
		repo: repo,
	}
}