// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"api-gateway/pkg/api"
	"api-gateway/pkg/api/handler"
	"api-gateway/pkg/client"
	"api-gateway/pkg/config"
)

// Injectors from wire.go:

func InitializeAPI(cfg config.Config) (*api.Server, error) {
	streamClient, err := client.NewStreamClient(cfg)
	if err != nil {
		return nil, err
	}
	streamHandler := handler.NewStreamHandler(streamClient)
	server := api.NewServerHTTP(cfg, streamHandler)
	return server, nil
}
