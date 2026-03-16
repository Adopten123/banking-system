package http

import (
	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	postService domain.PostService
}

func NewPostHandler(service domain.PostService) *PostHandler {
	return &PostHandler{
		postService: service,
	}
}

func (h *PostHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/posts", h.CreatePost)
}
