package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"stream-service/pkg/domain"
	"stream-service/pkg/file"
	"stream-service/pkg/models/request"
	repointerface "stream-service/pkg/repository/interfaces"
	"stream-service/pkg/usecase/interfaces"
	"time"

	"github.com/google/uuid"
)

type streamUseCase struct {
	repo        repointerface.StreamRepository
	fileHandler file.Handler
}

var (
	funcMaxWait = time.Second * 5
	uploadDir   = "./uploads/"
)

func NewStreamUseCase(repo repointerface.StreamRepository, fileHandler file.Handler) interfaces.StreamUseCase {
	return &streamUseCase{
		repo:        repo,
		fileHandler: fileHandler,
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

func (s *streamUseCase) UploadFileAsStream(ctx context.Context, fileID string,
	dataChan <-chan []byte, errChan chan error) {

	// create folder path and folder
	folderPath := generateFolderPath(fileID)
	if err := s.fileHandler.MkdirAll(folderPath, 0755); err != nil {
		errChan <- fmt.Errorf("failed to create directory for upload file: %w", err)
		return
	}

	// create file path and file
	filePath := generateFilePath(folderPath, fileID)
	file, err := s.fileHandler.Create(filePath)
	if err != nil {
		errChan <- fmt.Errorf("failed to create file details on server: %w", err)
		return
	}
	defer file.Close()

	// start reading data
	var buffer []byte
	for {
		select {
		case buffer = <-dataChan: // get data
			// write data on file
			_, err := file.Write(buffer)
			if err != nil {
				errChan <- fmt.Errorf("failed to write data on file: %w", err)
				return
			}
		case err = <-errChan:
			if err == io.EOF { // if EOF means stream completed so sending file id
				log.Println("EOF received on usecase and returning")
				return
			}
			log.Println("received error while receiving data on stream: ", err)
			return
		case <-ctx.Done(): // check return signal from context to cancel then return
			log.Println("usecase returned by cancel signal")
			return
		case <-time.After(funcMaxWait): // exit from the function if nothing happened for max function wait time
			log.Println("stream usecase function time out")
			return
		}
	}
}

// To generate folder path according to folder path and file id
func generateFolderPath(fileID string) string {
	return uploadDir + fileID
}

// To generate file path
func generateFilePath(folderPath, fileID string) string {
	return folderPath + "/" + fileID

}
