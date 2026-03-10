package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	_ "github.com/Adopten123/banking-system/service-account/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

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
		r.Post("/{id}/withdraw", h.withdraw)
		r.Post("/{id}/transfer", h.transfer)

		r.Post("/{id}/block", h.blockAccount)
		r.Post("/{id}/freeze", h.freezeAccount)
		r.Post("/{id}/activate", h.activateAccount)
		r.Post("/{id}/close", h.closeAccount)

		r.Put("/{id}/credit-limit", h.updateCreditLimit)

		r.Post("/{id}/cards", h.issueCard)
	})

	return r
}
