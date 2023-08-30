package gateway

import (
	"english/cmd/domain/repository"

	"gorm.io/gorm"
)

type ImgMySQLRepository struct {
	db *gorm.DB
}

func NewImgMySQLRepository(db *gorm.DB) repository.ImgRepository {
	return &ImgMySQLRepository{
		db: db,
	}
}
