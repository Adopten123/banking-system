package service

import "github.com/Adopten123/banking-system/service-community/internal/domain"

type profileService struct {
	repo domain.ProfileRepository
}

func NewProfileService(repo domain.ProfileRepository) domain.ProfileService {
	return &profileService{repo: repo}
}
