package tests

import (
	"net/http"

	transport_http "github.com/Adopten123/banking-system/service-account/internal/handler/http"

	"github.com/Adopten123/banking-system/service-account/internal/repository/postgres"
	"github.com/Adopten123/banking-system/service-account/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestEnv(dbPool *pgxpool.Pool) (http.Handler, *MockPublisher, *MockCardVaultClient) {
	// 1. Make mocks
	mockPublisher := NewMockPublisher()
	mockExchanger := &MockExchangeClient{}
	mockVault := &MockCardVaultClient{}
	// 2. Init repo with test DB pool
	repo := postgres.NewAccountRepo(dbPool)
	// 3. Starting service
	svc := service.NewAccountService(repo, mockPublisher, mockExchanger, mockVault)

	// 4. Init http layer
	handler := transport_http.NewHandler(svc)
	router := handler.InitRoutes()

	return router, mockPublisher, mockVault
}
