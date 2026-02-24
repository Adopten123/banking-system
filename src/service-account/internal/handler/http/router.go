package http

import (
	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/Adopten123/banking-system/service-account/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	service domain.AccountService
}

func NewHandler(service domain.AccountService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api/system", func(r chi.Router) {
		r.Get("/ping", h.ping)
	})

	r.Route("/api/accounts", func(r chi.Router) {
		r.Post("/", h.createAccount)

		r.Get("/{id}", h.getAccountBalance)
		r.Get("/{id}/transactions", h.getTransactions)

		r.Post("/{id}/deposit", h.deposit)
	})

	return r
}
