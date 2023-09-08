package main

import (
	"log"
	"stream-service/pkg/config"
	"stream-service/pkg/di"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	srv, err := di.InitializeAPI(cfg)
	if err != nil {
		log.Fatalf("failed to initialize api: %v", err)
	}

	if err = srv.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
