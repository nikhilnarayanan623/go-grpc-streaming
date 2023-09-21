//go:build wireinject
// +build wireinject

package di

import (
	"stream-service/pkg/api"
	"stream-service/pkg/api/service"
	"stream-service/pkg/config"
	"stream-service/pkg/db"
	"stream-service/pkg/file"
	"stream-service/pkg/repository"
	"stream-service/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*api.Server, error) {

	wire.Build(
		db.ConnectDatabase,
		repository.NewStreamRepository,
		file.NewHandler,
		usecase.NewStreamUseCase,
		service.NewStreamService,
		api.NewServerGRPC,
	)

	return &api.Server{}, nil
}
