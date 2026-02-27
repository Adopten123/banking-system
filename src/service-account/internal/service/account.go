package service

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

type AccountService struct {
	repo domain.AccountRepository
}

func NewAccountService(repo domain.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) CreateAccount(
	ctx context.Context,
	params domain.CreateAccountInput,
) (*domain.Account, error) {

	publicID := uuid.New()
	const defaultStatusID = 1

	acc := &domain.Account{
		PublicID:     publicID,
		UserID:       params.UserID,
		TypeID:       params.TypeID,
		StatusID:     defaultStatusID,
		CurrencyCode: params.CurrencyCode,
		Name:         params.Name,
	}

	createdAcc, err := s.repo.Create(ctx, acc)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return createdAcc, nil
}