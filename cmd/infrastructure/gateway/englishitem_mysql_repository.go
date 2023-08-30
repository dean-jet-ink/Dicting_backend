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
	return &EnglishItemMySQLRepository{
		db: db,
	}
}

func (er *EnglishItemMySQLRepository) Create(englishItem *model.EnglishItem) error {
	enti := er.modelToEntity(englishItem)

	if err := er.db.Create(enti).Error; err != nil {
		return err
	}

	return nil
}

func (er *EnglishItemMySQLRepository) modelToEntity(m *model.EnglishItem) *entity.EnglishItemEntity {

	exampleEntitys := []*entity.ExampleEntity{}
	for _, example := range m.Examples() {
		enti := &entity.ExampleEntity{
			Id:            example.Id,
			Example:       example.Example,
			Translation:   example.Translation,
			EnglishItemId: m.Id(),
		}
		exampleEntitys = append(exampleEntitys, enti)
	}

	imgEntitys := []*entity.ImgEntity{}
	for _, img := range m.Imgs() {
		enti := &entity.ImgEntity{
			Id:            img.Id(),
			URL:           img.URL(),
			EnglishItemId: m.Id(),
		}
		imgEntitys = append(imgEntitys, enti)
	}

	return &entity.EnglishItemEntity{
		Id:             m.Id(),
		Content:        m.Content(),
		JaTranslations: m.JoinJaTranslations(),
		EnExplanation:  m.EnExplanation(),
		UserId:         m.UserId(),
		Examples:       exampleEntitys,
		Imgs:           imgEntitys,
	}
}

func (er *EnglishItemMySQLRepository) entityToModel(e *entity.EnglishItemEntity, m *model.EnglishItem) {
	m.SetId(e.Id)
	m.SetContent(e.Content)
	m.SetJaTranslationsFromStr(e.JaTranslations)
	m.SetEnExplanation(e.EnExplanation)
	m.SetUserId(e.UserId)
}
