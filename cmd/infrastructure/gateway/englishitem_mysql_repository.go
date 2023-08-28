package gateway

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/infrastructure/entity"

	"gorm.io/gorm"
)

type EnglishItemMySQLRepository struct {
	db *gorm.DB
}

func NewEnglishItemMySQLReporitory(db *gorm.DB) repository.EnglishItemRepository {
	return &EnglishItemMySQLRepository{}
}

func (er *EnglishItemMySQLRepository) Create(englishItem *model.EnglishItem) error {
	enti := er.modelToEntity(englishItem)
	if err := er.db.Create(enti).Error; err != nil {
		return err
	}

	er.entityToModel(enti, englishItem)
	return nil
}

func (er *EnglishItemMySQLRepository) modelToEntity(m *model.EnglishItem) *entity.EnglishItemEntity {
	return entity.NewEnglishItemEntity(m.Id(), m.Content(), m.JoinJaTranslations(), m.EnExplanation(), m.UserId())
}

func (er *EnglishItemMySQLRepository) entityToModel(e *entity.EnglishItemEntity, m *model.EnglishItem) {
	m.SetId(e.Id)
	m.SetContent(e.Content)
	m.SetJaTranslationsFromStr(e.JaTranslations)
	m.SetEnExplanation(e.EnExplanation)
	m.SetUserId(e.UserId)
}
