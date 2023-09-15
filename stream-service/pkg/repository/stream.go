package repository

import (
	"context"
	"stream-service/pkg/domain"
	"stream-service/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type streamRepo struct {
	db *gorm.DB
}

func NewStreamRepository(db *gorm.DB) interfaces.StreamRepository {

	return &streamRepo{
		db: db,
	}
}

func (s *streamRepo) SaveFileDetails(ctx context.Context, details domain.FileDetails) error {

	query := `INSERT INTO file_details (id, name, content_type, uploaded_at) VALUES($1, $2, $3, $4)`
	return s.db.Exec(query, details.ID, details.Name, details.ContentType, details.UploadedAt).Error
}
