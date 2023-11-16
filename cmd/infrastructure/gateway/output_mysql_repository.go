package gateway

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/infrastructure/entity"
	"time"

	"gorm.io/gorm"
)

type OutputMySQLRepository struct {
	db *gorm.DB
}

func NewOutputRepository(db *gorm.DB) repository.OutputRepository {
	return &OutputMySQLRepository{
		db: db,
	}
}

func (r OutputMySQLRepository) FindOutputTimesByEnglishItemId(englishItemId string) ([]*time.Time, error) {
	entities := []*entity.OutputEntity{}
	if err := r.db.Model(&entity.OutputEntity{}).Where("english_item_id = ?", englishItemId).Group("created_at").Order("created_at DESC").Select("created_at").Find(&entities).Error; err != nil {
		return nil, err
	}

	times := []*time.Time{}

	for _, entity := range entities {
		times = append(times, &entity.CreatedAt)
	}

	return times, nil
}

func (r OutputMySQLRepository) FindByEnglishItemIdAndCreatedAt(englishItemId string, createdAt time.Time) ([]*model.Output, error) {
	entities := []*entity.OutputEntity{}
	if err := r.db.Where("english_item_id = ? AND created_at = ?", englishItemId, createdAt).Find(&entities).Error; err != nil {
		return nil, err
	}

	models := r.entitiesToModels(entities)

	return models, nil
}

func (r OutputMySQLRepository) Create(output *model.Output) error {
	entity := r.modelToEntity(output)

	if err := r.db.Create(entity).Error; err != nil {
		return err
	}

	return nil
}

func (r OutputMySQLRepository) Delete(englishItemId string, createdAt time.Time) error {
	if err := r.db.Where("english_item_id = ? AND created_at = ?", englishItemId, createdAt).Delete(&model.Output{}).Error; err != nil {
		return err
	}

	return nil
}

func (r OutputMySQLRepository) modelToEntity(output *model.Output) *entity.OutputEntity {
	return &entity.OutputEntity{
		Id:            output.ID(),
		Question:      output.Question(),
		Answer:        output.Answer(),
		Advice:        output.Advice(),
		Index:         output.Index(),
		EnglishItemId: output.EnglishItemId(),
		CreatedAt:     output.CreatedAt(),
	}
}

func (r OutputMySQLRepository) entitiesToModels(entities []*entity.OutputEntity) []*model.Output {
	models := []*model.Output{}

	for _, entity := range entities {
		model := model.NewOutput(entity.Id, entity.EnglishItemId, "", entity.Question, entity.Answer, entity.Advice, entity.Index, entity.CreatedAt)

		models = append(models, model)
	}

	return models
}
