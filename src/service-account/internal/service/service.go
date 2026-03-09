package service

import "github.com/Adopten123/banking-system/service-account/internal/domain"

type AccountService struct {
	repo      domain.AccountRepository
	publisher domain.EventPublisher
	exchanger domain.ExchangeRateClient
}

func NewAccountService(
	repo domain.AccountRepository,
	publisher domain.EventPublisher,
	exchanger domain.ExchangeRateClient,
) *AccountService {

	return &AccountService{
		repo:      repo,
		publisher: publisher,
		exchanger: exchanger,
	}
}
