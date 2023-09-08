package usecase

import (
	repointerface "stream-service/pkg/repository/interfaces"
	"stream-service/pkg/usecase/interfaces"
)

type streamUseCase struct {
	repo repointerface.StreamRepository
}

func NewStreamUseCase(repo repointerface.StreamRepository) interfaces.StreamUseCase {
	return &streamUseCase{
		repo: repo,
	}
}
