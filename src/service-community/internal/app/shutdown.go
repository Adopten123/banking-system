package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Operation func(ctx context.Context) error

func GracefulShutdown(timeout time.Duration, name string, ops ...Operation) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Printf("Graceful shutdown initiated for [%s]...", name)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, op := range ops {
		if err := op(ctx); err != nil {
			log.Printf("Error during shutdown operation: %v", err)
		}
	}

	log.Printf("[%s] shutdown complete. Exiting.", name)
}
