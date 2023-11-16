package repository

import (
	"english/cmd/domain/model"
	"time"
)

type OutputRepository interface {
	FindOutputTimesByEnglishItemId(englishItemId string) ([]*time.Time, error)
	FindByEnglishItemIdAndCreatedAt(englishItemId string, createdAt time.Time) ([]*model.Output, error)
	Create(output *model.Output) error
	Delete(englishItemId string, createdAt time.Time) error
}
