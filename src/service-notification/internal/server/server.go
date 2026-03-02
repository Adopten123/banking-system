package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown(cleanupFunc func()) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down service-notification gracefully...")

	if cleanupFunc != nil {
		cleanupFunc()
	}

	log.Println("Service stopped")
}
