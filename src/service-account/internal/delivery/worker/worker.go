package worker

import (
	"context"
	"log"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/repository/postgres"
	"github.com/Adopten123/banking-system/service-account/internal/service"
)

type RecurringWorker struct {
	repo    *postgres.AccountRepo
	service *service.AccountService
}

func NewRecurringWorker(repo *postgres.AccountRepo, svc *service.AccountService) *RecurringWorker {
	return &RecurringWorker{
		repo:    repo,
		service: svc,
	}
}

func (w *RecurringWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)

	log.Println("[WORKER] Recurring payments worker started")

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("[WORKER] Recurring payments worker stopped gracefully")
				return
			case <-ticker.C:
				w.ProcessDuePayments(ctx)
			}
		}
	}()
}
