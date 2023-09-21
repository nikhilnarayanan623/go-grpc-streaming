package domain

import (
	"time"

	"github.com/google/uuid"
)

type FileDetails struct {
	ID          uuid.UUID `gorm:"primaryKey;not null"`
	Name        string    `gorm:"not null"`
	ContentType string    `gorm:"not null"`
	UploadedAt  time.Time `gorm:"not null"`
}
