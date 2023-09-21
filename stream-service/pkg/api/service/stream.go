package service

import (
	"context"
	"io"
	"log"
	"stream-service/pkg/models/request"
	"stream-service/pkg/pb"
	"stream-service/pkg/usecase/interfaces"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *StreamService) Upload(stream pb.StreamService_UploadServer) error {

	// first take the file detail from the stream
	streamFile, err := stream.Recv()
	if err != nil {
		log.Println("failed to get file detail from stream")
		return status.Errorf(codes.InvalidArgument, "failed to receive file detail from stream: %v", err)
	}

	fileInfo := streamFile.GetInfo()
	if fileInfo == nil {
		return status.Errorf(codes.InvalidArgument, "provide file info on stream initially")
	}

	fileDetails := request.FileDetails{
		Name:        fileInfo.GetName(),
		ContentType: fileInfo.GetContentType(),
	}

	// create a context with cancel to send signal of closing to usecase
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// first upload the file details
	fileID, err := s.usecase.UploadFileDetails(ctx, fileDetails)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// run the usecase concurrently to read and upload file
	dataChan, errChan := make(chan []byte), make(chan error)

	go s.usecase.UploadFileAsStream(ctx, fileID, dataChan, errChan)

	// receive stream data
	for {
		select {
		// always check error on usecase
		case err = <-errChan:
			return status.Errorf(codes.Internal, "failed to store data: %v", err)
		default:
			// receive data from stream
			streamFile, err = stream.Recv()
			if err != nil {
				if err == io.EOF {

					log.Println("stream completed")
					// send EOF to notify stop waiting for data and send file id
					errChan <- io.EOF
					res := pb.UploadResponse{
						Id: fileID,
					}
					return stream.SendAndClose(&res)
				}
				// if error not EOF then return from stream
				return status.Errorf(codes.InvalidArgument, "failed to get stream file from client: %v", err)
			}
			// if no error to get stream
			// send data through data chan
			dataChan <- streamFile.GetData()
		}

	}
}
