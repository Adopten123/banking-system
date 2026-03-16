package postgres

import "github.com/Adopten123/banking-system/service-community/internal/domain"

type profileRepository struct {
	q *Queries
}

func NewProfileRepository(q *Queries) domain.ProfileRepository {
	return &profileRepository{q: q}
}