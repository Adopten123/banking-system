package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Adopten123/banking-system/service-community/internal/app"
	"github.com/Adopten123/banking-system/service-community/internal/config"
	transport "github.com/Adopten123/banking-system/service-community/internal/handler/http"
	"github.com/Adopten123/banking-system/service-community/internal/repository/postgres"
	"github.com/Adopten123/banking-system/service-community/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// @title           Community Service API
// @version         1.0
// @description     Социальное ядро экосистемы Banking System. Отвечает за посты, ленту и мессенджер.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8083
// @BasePath  /
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Fatal error loading config: %v", err)
	}

	log.Printf("Starting Community Service in [%s] mode...", cfg.Env)

	ctx := context.Background()

	// Init DB
	dbPool := app.InitPostgres(ctx, cfg.DB.URL)
	defer dbPool.Close()

	queries := postgres.New(dbPool)

	postRepo := postgres.NewPostRepository(queries)
	postService := service.NewPostService(postRepo)
	postHandler := transport.NewPostHandler(postService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	postHandler.RegisterRoutes(r)

	// HTTP server
	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Run server
	go func() {
		log.Printf("HTTP server is listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful Shutdown
	app.GracefulShutdown(
		5*time.Second,
		"Community Service",
		server.Shutdown,
	)
}
