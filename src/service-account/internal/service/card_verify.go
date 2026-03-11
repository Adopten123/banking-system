package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

func (s *AccountService) VerifyCardForPayment(ctx context.Context, input domain.VerifyCardInput) (*domain.VerifyCardResult, error) {
	// 1. Checking valid in vault
	isValid, tokenIDStr, err := s.vault.VerifyCard(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("vault verification failed: %w", err)
	}

	if !isValid || tokenIDStr == "" {
		return &domain.VerifyCardResult{IsValid: false}, nil
	}

	// 2. Parse token
	cardID, err := uuid.Parse(tokenIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid token received from vault: %w", err)
	}

	// 3. Getting card in core
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		if errors.Is(err, domain.ErrCardNotFound) {
			return &domain.VerifyCardResult{IsValid: false}, nil
		}
		return nil, err
	}

	// 4. Checking status
	if card.Status != "active" {
		return &domain.VerifyCardResult{IsValid: false}, nil
	}

	return &domain.VerifyCardResult{
		IsValid:   true,
		CardID:    card.ID,
		AccountID: card.AccountID,
	}, nil
}
