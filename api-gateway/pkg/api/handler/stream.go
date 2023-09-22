package handler

import (
	"api-gateway/pkg/api/handler/interfaces"
	clientinterface "api-gateway/pkg/client/interfaces"
	"api-gateway/pkg/models/request"
	"api-gateway/pkg/utils"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type streamHandler struct {
	client clientinterface.StreamClient
}

func NewStreamHandler(client clientinterface.StreamClient) interfaces.StreamHandler {
	return &streamHandler{
		client: client,
	}
}

func (s *streamHandler) Upload(ctx echo.Context) error {
	// get file name from form values
	name := ctx.FormValue("name")
	if name == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": "File name not provided",
		})
	}
	// get file as form file
	fh, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{

			"message": "Failed to get file from request",
			"error":   err.Error(),
		})
	}
	// get the file content type
	contentType, err := utils.GetContentType(fh)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{

			"message": "Failed to get file content type",
			"error":   err.Error(),
		})
	}

	fileDetails := request.FileDetails{
		Name:        name,
		FileHeader:  fh,
		ContentType: contentType,
	}
	/// upload the file to client
	id, err := s.client.Upload(context.Background(), fileDetails)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{

			"message": "Failed  upload file",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{

		"message": "File upload completed",
		"File ID": id,
	})

}
