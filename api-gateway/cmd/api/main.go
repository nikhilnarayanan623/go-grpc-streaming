package main

import (
	"api-gateway/pkg/api"
	"api-gateway/pkg/config"
	"log"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	server := api.NewServerHTTP(cfg)

	server.Start()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
