package repository

import (
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
