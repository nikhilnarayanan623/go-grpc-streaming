package usecase

import (
	"context"
	"fmt"
	"stream-service/pkg/domain"
	"stream-service/pkg/models/request"
	repointerface "stream-service/pkg/repository/interfaces"
	"stream-service/pkg/usecase/interfaces"
	"time"

	"github.com/google/uuid"
)

type streamUseCase struct {
	repo repointerface.StreamRepository
}

func NewStreamUseCase(repo repointerface.StreamRepository) interfaces.StreamUseCase {
	return &streamUseCase{
		repo: repo,
	}
}

func (s *streamUseCase) UploadFileDetails(ctx context.Context, details request.FileDetails) (string, error) {

	fileID := uuid.New()

	fileDetails := domain.FileDetails{
		ID:          fileID,
		Name:        details.Name,
		ContentType: details.ContentType,
		UploadedAt:  time.Now(),
	}

	// save file details on database
	err := s.repo.SaveFileDetails(ctx, fileDetails)
	if err != nil {
		return "", fmt.Errorf("failed to save file details on database: \n%w", err)
	}

	return fileID.String(), nil
}
