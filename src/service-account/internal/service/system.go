package service

import (
	"context"
)

func (s *AccountService) CheckHealth(ctx context.Context) string {
	err := s.repo.Ping(ctx)
	if err != nil {
		return "error"
	}
	return "database and service are ok"
}
