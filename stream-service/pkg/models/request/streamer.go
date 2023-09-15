package request

type FileDetails struct {
	Name        string `validator:"required,min=3"`
	ContentType string `validator:"required"`
}
