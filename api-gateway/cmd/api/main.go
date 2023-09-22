package main

import (
	"api-gateway/pkg/config"
	"api-gateway/pkg/di"
	"log"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	server, err := di.InitializeAPI(cfg)
	if err != nil {
		log.Fatalf("failed to initialize api : %v", err)
	}

	server.Start()
}
