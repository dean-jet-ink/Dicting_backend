package gateway

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/infrastructure/entity"
	"fmt"

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

func (er *EnglishItemMySQLRepository) FindByUserIdAndContent(userId, content string) ([]*model.EnglishItem, error) {
	entities := []*entity.EnglishItemEntity{}
	if err := er.db.Preload("Imgs").Preload("Examples").Preload("Proficiencies").Where("content = ?", fmt.Sprintf("%v%%", content)).Find(&entities).Error; err != nil {
		return nil, err
	}

	englishItems := []*model.EnglishItem{}
	for _, entity := range entities {
		englishItem := &model.EnglishItem{}
		er.entityToModel(entity, englishItem)
		englishItems = append(englishItems, englishItem)
	}

	return englishItems, nil
}

func (er *EnglishItemMySQLRepository) FindByUserId() ([]*model.EnglishItem, error) {
	entities := []*entity.EnglishItemEntity{}
	if err := er.db.Preload("Imgs").Preload("Proficiencies").Where("isThumbnail = ?", true).Find(&entities).Error; err != nil {
		return nil, err
	}

	englishItems := []*model.EnglishItem{}
	for _, entity := range entities {
		englishItem := &model.EnglishItem{}
		er.entityToModel(entity, englishItem)
		englishItems = append(englishItems, englishItem)
	}

	return englishItems, nil
}

func (er *EnglishItemMySQLRepository) modelToEntity(m *model.EnglishItem) *entity.EnglishItemEntity {
	exampleEntities := []*entity.ExampleEntity{}
	for _, example := range m.Examples() {
		enti := &entity.ExampleEntity{
			Id:            example.Id,
			Example:       example.Example,
			Translation:   example.Translation,
			EnglishItemId: m.Id(),
		}
		exampleEntities = append(exampleEntities, enti)
	}

	imgEntities := []*entity.ImgEntity{}
	for _, img := range m.Imgs() {
		enti := &entity.ImgEntity{
			Id:            img.Id(),
			URL:           img.URL(),
			EnglishItemId: m.Id(),
		}
		imgEntities = append(imgEntities, enti)
	}

	return &entity.EnglishItemEntity{
		Id:            m.Id(),
		Content:       m.Content(),
		Translations:  m.JoinTranslations(),
		EnExplanation: m.EnExplanation(),
		UserId:        m.UserId(),
		Examples:      exampleEntities,
		Imgs:          imgEntities,
	}
}

func (er *EnglishItemMySQLRepository) entityToModel(e *entity.EnglishItemEntity, m *model.EnglishItem) {
	m.SetId(e.Id)
	m.SetContent(e.Content)
	m.SetTranslationsFromStr(e.Translations)
	m.SetEnExplanation(e.EnExplanation)
	m.SetUserId(e.UserId)
	m.SetProficiency(model.Proficiency(e.Proficiency.Proficiency))
	m.SetExp(e.Proficiency.Exp)

	exampleDomains := []*model.Example{}
	for _, example := range e.Examples {
		exampleDomain := model.NewExample(example.Id, example.Example, example.Translation)
		exampleDomains = append(exampleDomains, exampleDomain)
	}
	m.SetExamples(exampleDomains)
	imgDomains := []*model.Img{}
	for _, img := range e.Imgs {
		imgDomain := model.NewImg(img.Id, img.URL, img.IsThumbnail)
		imgDomains = append(imgDomains, imgDomain)
	}
	m.SetImgs(imgDomains)
}
