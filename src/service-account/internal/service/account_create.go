package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

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

	err = s.publisher.PublishAccountCreated(ctx, domain.AccountCreatedEvent{
		AccountID: createdAcc.ID,
		PublicID:  createdAcc.PublicID,
		UserID:    createdAcc.UserID,
		Currency:  createdAcc.CurrencyCode,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		fmt.Printf("ERROR: Failed to publish AccountCreated event: %v\n", err)
	}

	return createdAcc, nil
}
