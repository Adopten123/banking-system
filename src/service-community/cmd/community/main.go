package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Adopten123/banking-system/service-community/internal/app"
	"github.com/Adopten123/banking-system/service-community/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Fatal error loading config: %v", err)
	}

	log.Printf("Starting Community Service in [%s] mode...", cfg.Env)

	// TODO: init db, redis, rabbit

	// TODO: init layers

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// HTTP server
	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
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
