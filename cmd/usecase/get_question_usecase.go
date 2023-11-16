package usecase

import (
	"english/cmd/domain/api"
	"english/cmd/usecase/dto"
)

type GetQuestionUsecase interface {
	GetQuestion(input *dto.GetQuestionInput) (*dto.GetQuestionOutput, error)
}

type GetQuestionUsecaseImpl struct {
	chatAIAPI api.ChatAIAPI
}

func NewGetQuestionUsecase(chatAIAPI api.ChatAIAPI) GetQuestionUsecase {
	return &GetQuestionUsecaseImpl{
		chatAIAPI: chatAIAPI,
	}
}

func (u GetQuestionUsecaseImpl) GetQuestion(input *dto.GetQuestionInput) (*dto.GetQuestionOutput, error) {
	res, err := u.chatAIAPI.GetQuestion(input.Content)
	if err != nil {
		return nil, err
	}

	output := &dto.GetQuestionOutput{
		Question: res,
	}

	return output, nil
}
