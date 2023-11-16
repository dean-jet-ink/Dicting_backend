package usecase

import (
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
)

type GetOutputTimesUsecase interface {
	GetOutputTimes(englishItemId string) (*dto.GetOutputTimesOutput, error)
}

type GetOutputTimesUsecaseImpl struct {
	outputRepo repository.OutputRepository
}

func NewGetOutputTimesUsecase(outputRepo repository.OutputRepository) GetOutputTimesUsecase {
	return &GetOutputTimesUsecaseImpl{
		outputRepo: outputRepo,
	}
}

func (u *GetOutputTimesUsecaseImpl) GetOutputTimes(englishItemId string) (*dto.GetOutputTimesOutput, error) {
	times, err := u.outputRepo.FindOutputTimesByEnglishItemId(englishItemId)
	if err != nil {
		return nil, err
	}

	output := &dto.GetOutputTimesOutput{
		OutputTimes: times,
	}

	return output, nil
}
