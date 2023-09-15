package interfaces

import (
	"context"
	"stream-service/pkg/models/request"
)

type StreamUseCase interface {
	UploadFileDetails(ctx context.Context, details request.FileDetails) (string, error)
}
