package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Adopten123/banking-system/service-community/internal/config"
	deliveryRMQ "github.com/Adopten123/banking-system/service-community/internal/delivery/rabbitmq"
	deliveryWS "github.com/Adopten123/banking-system/service-community/internal/delivery/websocket"
	transport "github.com/Adopten123/banking-system/service-community/internal/handler/http"
	"github.com/Adopten123/banking-system/service-community/internal/repository/postgres"
	"github.com/Adopten123/banking-system/service-community/internal/service"
)

func Run(cfg *config.Config) {
	ctx := context.Background()

	// 1. Init DB
	dbPool := InitPostgres(ctx, cfg.DB.URL)
	defer dbPool.Close()

	queries := postgres.New(dbPool)

	// 2. Dependency Injection
	postRepo := postgres.NewPostRepository(queries)
	postService := service.NewPostService(postRepo)
	postHandler := transport.NewPostHandler(postService)

	profileRepo := postgres.NewProfileRepository(queries)
	profileService := service.NewProfileService(profileRepo)

	chatRepo := postgres.NewChatRepository(queries)
	chatService := service.NewChatService(chatRepo)
	chatHandler := transport.NewChatHandler(chatService)

	ssoConsumer, err := deliveryRMQ.NewSSOConsumer(cfg.RabbitMQ.URL, profileService)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}
	ssoConsumer.Start(ctx)

	wsConfig := config.WebSocketConfig{
		MaxMessageSize: 4096,
		PongWait:       60 * time.Second,
		WriteWait:      10 * time.Second,
	}

	hub := deliveryWS.NewHub(wsConfig, chatService)
	go hub.Run()
	wsHandler := deliveryWS.NewWSHandler(hub)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	setupRoutes(r, wsHandler, postHandler, chatHandler)

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("HTTP server is listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	GracefulShutdown(
		5*time.Second,
		"Community Service",
		server.Shutdown,
	)
}

func setupRoutes(
	r chi.Router,
	wsHandler *deliveryWS.WSHandler,
	postHandler *transport.PostHandler,
	chatHandler *transport.ChatHandler,
) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Get("/api/v1/ws", wsHandler.ServeWS)

	postHandler.RegisterRoutes(r)
	chatHandler.RegisterRoutes(r)
}
