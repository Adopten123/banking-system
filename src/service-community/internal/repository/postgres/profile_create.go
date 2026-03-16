package postgres

import (
	"context"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *profileRepository) Create(ctx context.Context, input domain.CreateProfileInput) (*domain.Profile, error) {
	params := CreateUserProfileParams{
		UserID:      pgtype.UUID{Bytes: input.UserID, Valid: true},
		Username:    input.Username,
		IsVerified:  false,
		IsStaff:     false,
		IsPrivate:   false,
	}

	if input.DisplayName != "" {
		params.DisplayName = pgtype.Text{String: input.DisplayName, Valid: true}
	}

	row, err := r.q.CreateUserProfile(ctx, params)
	if err != nil {
		return nil, err
	}

	var dispName *string
	if row.DisplayName.Valid {
		dispName = &row.DisplayName.String
	}

	return &domain.Profile{
		UserID:      row.UserID.Bytes,
		Username:    row.Username,
		DisplayName: dispName,
		IsVerified:  row.IsVerified,
		IsStaff:     row.IsStaff,
		IsPrivate:   row.IsPrivate,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}, nil
}