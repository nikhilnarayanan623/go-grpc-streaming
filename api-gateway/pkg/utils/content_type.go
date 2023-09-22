package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func GetContentType(fileHeader *multipart.FileHeader) (string, error) {

	if fileHeader == nil {
		return "", errors.New("got nil file header")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	// first 512 bytes needed for http detectContentType function to detect content type
	minSizeForDetectContent := 512

	buffer := make([]byte, minSizeForDetectContent)
	_, err = io.ReadAtLeast(file, buffer, minSizeForDetectContent)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			return "", fmt.Errorf("file size is too low to detect content type")
		}
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
