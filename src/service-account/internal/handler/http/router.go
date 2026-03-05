package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/Adopten123/banking-system/service-account/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

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

		r.Get("/{id}", h.getAccountInfo)
		r.Get("/{id}/transactions", h.getTransactions)
		r.Get("/{id}/balance", h.getAccountBalance)

		r.Post("/{id}/deposit", h.deposit)
		r.Post("/{id}/transfer", h.transfer)

		r.Post("/{id}/block", h.blockAccount)
		r.Post("/{id}/close", h.closeAccount)

		r.Put("/{id}/credit-limit", h.updateCreditLimit)
	})

	return r
}
