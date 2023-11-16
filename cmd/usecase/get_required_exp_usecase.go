package usecase

import (
	"english/cmd/domain/model"
	"english/cmd/usecase/dto"
)

type GetRequiredExpUsecase interface {
	GetRequiredExp() *dto.GetRequiredExpOutput
}

type GetRequiredExpUsecaseImpl struct {
}

func NewGetRequiredExpUsecase() GetRequiredExpUsecase {
	return &GetRequiredExpUsecaseImpl{}
}

func (u *GetRequiredExpUsecaseImpl) GetRequiredExp() *dto.GetRequiredExpOutput {
	output := &dto.GetRequiredExpOutput{
		UnderstandExp: model.UnderstandExp,
		MasteredExp:   model.MasteredExp,
	}

	return output
}
