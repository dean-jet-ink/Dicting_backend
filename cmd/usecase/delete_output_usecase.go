package usecase

import (
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
)

type DeleteOutputUsecase interface {
	Delete(input *dto.DeleteOutputInput) error
}

type DeleteOutputUsecaseImpl struct {
	outputRepo repository.OutputRepository
}

func NewDeleteOutputUsecase(outputRepo repository.OutputRepository) DeleteOutputUsecase {
	return &DeleteOutputUsecaseImpl{
		outputRepo: outputRepo,
	}
}

func (u *DeleteOutputUsecaseImpl) Delete(input *dto.DeleteOutputInput) error {
	if err := u.outputRepo.Delete(input.EnglishItemId, input.CreatedAt); err != nil {
		return err
	}

	return nil
}
