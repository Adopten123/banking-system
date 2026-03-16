package postgres

import "github.com/Adopten123/banking-system/service-community/internal/domain"

type postRepository struct {
	q *Queries
}

func NewPostRepository(q *Queries) domain.PostRepository {
	return &postRepository{
		q: q,
	}
}