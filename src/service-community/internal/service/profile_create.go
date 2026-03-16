package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/google/uuid"
)

func (s *profileService) CreateProfileFromEvent(ctx context.Context, event domain.UserCreatedEvent) error {
	userUUID, err := uuid.Parse(event.UserID)
	if err != nil {
		return fmt.Errorf("invalid user uuid in event: %w", err)
	}

	input := domain.CreateProfileInput{
		UserID:      userUUID,
		Username:    event.Username,
		DisplayName: event.DisplayName,
	}

	_, err = s.repo.Create(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	log.Printf("Successfully created profile for user: %s", event.Username)
	return nil
}
