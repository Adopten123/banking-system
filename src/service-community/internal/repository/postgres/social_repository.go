package postgres

import "github.com/Adopten123/banking-system/service-community/internal/domain"

type socialRepository struct {
	q *Queries
}

func NewSocialRepository(q *Queries) domain.SocialRepository {
	return &socialRepository{q: q}
}