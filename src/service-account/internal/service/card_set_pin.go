package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

var pinRegex = regexp.MustCompile(`^\d{4}$`)

func (s *AccountService) SetCardPin(ctx context.Context, cardID uuid.UUID, pin string) error {
	// 1. Validate
	if !pinRegex.MatchString(pin) {
		return domain.ErrInvalidPINFormat
	}

	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return err
	}

	if card.Status != "active" {
		return domain.ErrCardBlocked
	}

	// 2. Send PIN to vault
	err = s.vault.SetPin(ctx, cardID.String(), pin)
	if err != nil {
		return fmt.Errorf("failed to set secure PIN in vault: %w", err)
	}

	return nil
}