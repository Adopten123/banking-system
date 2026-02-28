package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
)

func (s *AccountService) UpdateCreditLimit(ctx context.Context, publicID uuid.UUID, limitStr string) error {
	limitFloat, err := strconv.ParseFloat(limitStr, 64)
	if err != nil || limitFloat < 0 {
		return fmt.Errorf("invalid credit limit: must be a non-negative number")
	}

	acc, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	if acc.StatusID != 1 {
		return domain.ErrAccountInactive
	}

	return s.repo.UpdateCreditLimit(ctx, acc.ID, limitStr)
}
