package gateway

import (
	"english/cmd/domain/repository"

	"gorm.io/gorm"
)

type ExampleMySQLRepository struct {
	db *gorm.DB
}

func NewExampleMySQLRepository(db *gorm.DB) repository.ExampleRepository {
	return &ExampleMySQLRepository{
		db: db,
	}
}
