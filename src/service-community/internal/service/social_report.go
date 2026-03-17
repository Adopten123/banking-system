package service

import (
	"context"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
)

func (s *socialService) FileReport(ctx context.Context, input domain.CreateReportInput) (*domain.Report, error) {
	return s.repo.CreateReport(ctx, input)
}