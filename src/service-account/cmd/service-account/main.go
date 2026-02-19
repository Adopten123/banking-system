package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.Info("Starting Service Account")

	fmt.Println("Server is running...")
}
