package domain

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	UserID      uuid.UUID
	Username    string
	DisplayName *string
	AvatarURL   *string
	Bio         *string
	IsVerified  bool
	IsStaff     bool
	IsPrivate   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateProfileInput struct {
	UserID      uuid.UUID
	Username    string
	DisplayName string
}
