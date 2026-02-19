package main

import (
	"log"
	"net/http"

	transport "github.com/Adopten123/banking-system/service-account/internal/handler/http"
	"github.com/Adopten123/banking-system/service-account/internal/repository/postgres"
	"github.com/Adopten123/banking-system/service-account/internal/service"
)

func main() {
	log.Println("Initializing service-account components...")

	repo := postgres.NewAccountRepo()

	svc := service.NewAccountService(repo)

	handler := transport.NewHandler(svc)
	router := handler.InitRoutes()

	log.Println("Server is listening on port :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}
