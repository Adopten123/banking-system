package service

import "github.com/Adopten123/banking-system/service-community/internal/domain"

type socialService struct {
	repo domain.SocialRepository
}

func NewSocialService(repo domain.SocialRepository) domain.SocialService {
	return &socialService{repo: repo}
}