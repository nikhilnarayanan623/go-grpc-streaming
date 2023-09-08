package service

import (
	"stream-service/pkg/pb"
	"stream-service/pkg/usecase/interfaces"
)

type StreamService struct {
	pb.UnimplementedStreamServiceServer
	usecase interfaces.StreamUseCase
}

func NewStreamService(usecase interfaces.StreamUseCase) pb.StreamServiceServer {
	return &StreamService{
		usecase: usecase,
	}
}
