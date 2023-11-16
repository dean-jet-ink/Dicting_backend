package usecase

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
	"english/lib"
	"time"
)

type CreateOutputUsecase interface {
	Create(input *dto.CreateOutputInput) error
}

type CreateOutputUsecaseImpl struct {
	outputRepo repository.OutputRepository
}

func NewCreateOutputUsecase(outputRepo repository.OutputRepository) CreateOutputUsecase {
	return &CreateOutputUsecaseImpl{
		outputRepo: outputRepo,
	}
}

func (u CreateOutputUsecaseImpl) Create(input *dto.CreateOutputInput) error {
	createdAt := time.Now()

	for _, output := range input.Outputs {
		guid, err := lib.GenerateULID()
		if err != nil {
			return err
		}

		outputModel := model.NewOutput(
			guid,
			input.EnglishItemId,
			"",
			output.Question,
			output.Answer,
			output.Advice,
			output.Index,
			createdAt,
		)

		if err := u.outputRepo.Create(outputModel); err != nil {
			return err
		}
	}

	return nil
}
