package request

import "mime/multipart"

type FileDetails struct {
	Name        string `validator:"required,min=3"`
	ContentType string `validator:"required"`
	FileHeader  *multipart.FileHeader
}
