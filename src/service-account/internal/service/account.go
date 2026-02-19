package service

import "github.com/Adopten123/banking-system/service-account/internal/domain"

type AccountService struct {
	repo domain.AccountRepository
}

func NewAccountService(repo domain.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) CheckHealth() string {
	err := s.repo.Ping()
	if err != nil {
		return "error"
	}
	return "database and service are ok"
}