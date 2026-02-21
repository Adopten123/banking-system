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

func (s *AccountService) CheckHealth(ctx context.Context) string {
	err := s.repo.Ping(ctx)
	if err != nil {
		return "error"
	}
	return "database and service are ok"
}

func (s *AccountService) CreateAccount(
	ctx context.Context,
	userID uuid.UUID,
	typeID int32,
	currencyCode, name string,
) (*domain.Account, error) {

	publicID := uuid.New()
	const defaultStatusID = 1

	acc := &domain.Account{
		PublicID:     publicID,
		UserID:       userID,
		TypeID:       typeID,
		StatusID:     defaultStatusID,
		CurrencyCode: currencyCode,
		Name:         name,
	}

	createdAcc, err := s.repo.Create(ctx, acc)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return createdAcc, nil
}
