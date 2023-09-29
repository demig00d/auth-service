package main

import (
	"log"

	"github.com/demig00d/auth-service/config"
	"github.com/demig00d/auth-service/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
