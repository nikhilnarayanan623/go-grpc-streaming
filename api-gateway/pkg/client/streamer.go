package client

import (
	"api-gateway/pkg/client/interfaces"
	"api-gateway/pkg/config"
	"api-gateway/pkg/models/request"
	"api-gateway/pkg/pb"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type streamClient struct {
	client pb.StreamServiceClient
}

var streamSize = 500

func NewStreamClient(cfg config.Config) (interfaces.StreamClient, error) {

	addr := fmt.Sprintf("%s:%s", cfg.StreamServiceHost, cfg.StreamServicePort)

	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial to stream service: %w", err)
	}
	client := pb.NewStreamServiceClient(cc)

	return &streamClient{
		client: client,
	}, nil
}

func (c *streamClient) Upload(ctx context.Context, fileDetails request.FileDetails) (string, error) {

	// get the stream service
	streamSvc, err := c.client.Upload(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to call upload method for stream client: %w", err)
	}

	fileRequest := pb.UploadRequest_Info{

		Info: &pb.FileMetaData{
			Name:        fileDetails.Name,
			ContentType: fileDetails.ContentType,
		},
	}

	// first send file meta data
	err = streamSvc.Send(&pb.UploadRequest{
		File: &fileRequest,
	})
	if err != nil {
		return "", fmt.Errorf("failed to send file details: %w", err)
	}

	file, err := fileDetails.FileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	data := make([]byte, streamSize)
	for {
		// read file data
		_, err := file.Read(data)
		if err != nil {
			// if file read completed break
			if err == io.EOF {
				log.Println("file read completed and stop streaming..")
				break
			}
			return "", fmt.Errorf("failed to read from file: %w", err)
		}
		// send stream data
		streamData := pb.UploadRequest{
			File: &pb.UploadRequest_Data{Data: data},
		}
		err = streamSvc.Send(&streamData)
		if err != nil {
			return "", fmt.Errorf("failed to send stream to server: %w", err)
		}
	}

	// close streaming
	res, err := streamSvc.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("failed to close streaming: %w", err)
	}

	return res.GetId(), nil
}
