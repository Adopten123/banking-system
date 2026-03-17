package http

import (
	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/go-chi/chi/v5"
)

type SocialHandler struct {
	service domain.SocialService
}

func NewSocialHandler(service domain.SocialService) *SocialHandler {
	return &SocialHandler{service: service}
}

func (h *SocialHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/users/{id}/follow", h.follow)
	r.Delete("/api/v1/users/{id}/follow", h.unfollow)
	r.Post("/api/v1/posts/{id}/like", h.like)
	r.Delete("/api/v1/posts/{id}/like", h.unlike)
	r.Post("/api/v1/reports", h.report)
}