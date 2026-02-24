package main

import (
	"context"
	"log"
	"os"

	"github.com/Adopten123/banking-system/service-account/internal/config"
	transport "github.com/Adopten123/banking-system/service-account/internal/handler/http"
	"github.com/Adopten123/banking-system/service-account/internal/repository/postgres"
	"github.com/Adopten123/banking-system/service-account/internal/server"
	"github.com/Adopten123/banking-system/service-account/internal/service"
)

const defaultCfgPath = "config/local.yaml"

// @title Banking System - Account Service API
// @version 1.0
// @description Микросервис для управления счетами и балансами.
// @host localhost:8081
// @BasePath /
func main() {
	// Init cfg path
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = defaultCfgPath
	}

	// Load cfg
	cfg := config.Load(cfgPath)
	log.Println("Config loaded successfully")

	// Init ctx
	ctx := context.Background()

	// Connect to DB
	pool, err := postgres.NewDBPool(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize DB pool: %v", err)
	}
	defer pool.Close()

	// Init layers
	repo := postgres.NewAccountRepo(pool)
	svc := service.NewAccountService(repo)
	handler := transport.NewHandler(svc)
	router := handler.InitRoutes()

	srv := server.NewServer(":8080", router)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server shutdown with error: %v", err)
	}

	log.Println("Server stopped properly")
}
