//go:build wireinject
// +build wireinject

package di

import (
	"api-gateway/pkg/api"
	"api-gateway/pkg/api/handler"
	"api-gateway/pkg/client"
	"api-gateway/pkg/config"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*api.Server, error) {

	wire.Build(
		client.NewStreamClient,
		handler.NewStreamHandler,
		api.NewServerHTTP,
	)

	return &api.Server{}, nil
}
