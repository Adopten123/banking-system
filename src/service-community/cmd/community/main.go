package main

import (
	"log"

	_ "github.com/Adopten123/banking-system/service-community/docs"
	"github.com/Adopten123/banking-system/service-community/internal/app"
	"github.com/Adopten123/banking-system/service-community/internal/config"
)

// @title           Community Service API
// @version         1.0
// @description     Social Network of Banking System
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8083
// @BasePath  /
func main() {
	// Read cfg
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Fatal error loading config: %v", err)
	}

	log.Printf("Starting Community Service in [%s] mode...", cfg.Env)

	// Run service
	app.Run(cfg)
}