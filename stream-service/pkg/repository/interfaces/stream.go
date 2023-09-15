package interfaces

import (
	"context"
	"stream-service/pkg/domain"
)

type StreamRepository interface {
	SaveFileDetails(ctx context.Context, details domain.FileDetails) error
}
