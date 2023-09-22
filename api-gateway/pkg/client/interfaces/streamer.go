package interfaces

import (
	"api-gateway/pkg/models/request"
	"context"
)

type StreamClient interface {
	Upload(ctx context.Context, file request.FileDetails) (string, error)
}
