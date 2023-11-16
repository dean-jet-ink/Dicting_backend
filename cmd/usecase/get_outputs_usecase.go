package usecase

import (
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
)

type GetOutputsUsecase interface {
	GetOutputs(input *dto.GetOutputsInput) (*dto.GetOutputsOutput, error)
}

type GetOutputsUsecaseImpl struct {
	outputRepo repository.OutputRepository
}

func NewGetOutputsUsecase(outputRepo repository.OutputRepository) GetOutputsUsecase {
	return &GetOutputsUsecaseImpl{
		outputRepo: outputRepo,
	}
}

func (u *GetOutputsUsecaseImpl) GetOutputs(input *dto.GetOutputsInput) (*dto.GetOutputsOutput, error) {
	outputs, err := u.outputRepo.FindByEnglishItemIdAndCreatedAt(input.EnglishItemId, input.CreatedAt)
	if err != nil {
		return nil, err
	}

	outputsDTO := &dto.GetOutputsOutput{}

	for _, output := range outputs {
		outputDTO := &dto.Output{
			Id:        output.ID(),
			Index:     output.Index(),
			Question:  output.Question(),
			Answer:    output.Answer(),
			Advice:    output.Advice(),
			CreatedAt: output.CreatedAt(),
		}

		outputsDTO.Outputs = append(outputsDTO.Outputs, outputDTO)
	}

	return outputsDTO, nil
}
