package http

import "github.com/Adopten123/banking-system/service-account/internal/domain"

type Handler struct {
	service domain.AccountService
}

func NewHandler(service domain.AccountService) *Handler {
	return &Handler{
		service: service,
	}
}
