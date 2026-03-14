package service

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

// resolveAccount - find account by source
func (s *AccountService) resolveAccount(
	ctx context.Context,
	entityType string,
	entityValue string,
) (*domain.Account, int32, uuid.UUID, error) {

	switch entityType {
	case "account":
		accUUID, err := uuid.Parse(entityValue)
		if err != nil {
			return nil, 0, uuid.Nil, domain.ErrInvalidFormat
		}

		acc, err := s.repo.GetByPublicID(ctx, accUUID)
		if err != nil {
			return nil, 0, uuid.Nil, err
		}
		if acc.StatusID != 1 {
			return nil, 0, uuid.Nil, domain.ErrAccountInactive
		}
		return acc, domain.SourceTypeIDAccount, accUUID, nil

	case "card":
		cardTokenStr, err := s.vault.GetTokenByPan(ctx, entityValue)
		if err != nil {
			return nil, 0, uuid.Nil, fmt.Errorf("card not found in vault: %w", err)
		}
		cardUUID, err := uuid.Parse(cardTokenStr)
		if err != nil {
			return nil, 0, uuid.Nil, fmt.Errorf("vault returned invalid token: %w", err)
		}
		card, err := s.repo.GetCardByID(ctx, cardUUID)
		if err != nil {
			return nil, 0, uuid.Nil, err
		}
		if card.Status != "active" {
			return nil, 0, uuid.Nil, domain.ErrCardBlocked
		}
		acc, err := s.repo.GetAccountInternalByID(ctx, card.AccountID)
		if err != nil {
			return nil, 0, uuid.Nil, fmt.Errorf("failed to get account for card: %w", err)
		}

		return acc, domain.SourceTypeIDCard, cardUUID, nil

	default:
		return nil, 0, uuid.Nil, fmt.Errorf("unsupported entity type: %s", entityType)
	}
}
