package domain

import (
	"context"
)

type ProfileRepository interface {
	Create(ctx context.Context, input CreateProfileInput) (*Profile, error)
}

type ProfileService interface {
	CreateProfileFromEvent(ctx context.Context, event UserCreatedEvent) error
}