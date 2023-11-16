package usecase

import (
	"english/cmd/domain/api"
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
	"time"
)

type AnswerQuestionsUsecase interface {
	AnswerQuestions(input *dto.AnswerQuestionsInput) (*dto.AnswerQuestionsOutput, error)
}

type AnswerQuestionsUsecaseImpl struct {
	englishItemRepo repository.EnglishItemRepository
	chatAIAPI       api.ChatAIAPI
}

func NewAnswerQuestionsUsecase(englishItemRepo repository.EnglishItemRepository, chatAIAPI api.ChatAIAPI) AnswerQuestionsUsecase {
	return &AnswerQuestionsUsecaseImpl{
		englishItemRepo: englishItemRepo,
		chatAIAPI:       chatAIAPI,
	}
}

func (u AnswerQuestionsUsecaseImpl) AnswerQuestions(input *dto.AnswerQuestionsInput) (*dto.AnswerQuestionsOutput, error) {
	englishItemId, content, answers := input.EnglishItemId, input.Content, input.Answers

	// Expの計算、更新
	englishItem, err := u.englishItemRepo.FindById(englishItemId)
	if err != nil {
		return nil, err
	}

	exp := englishItem.Exp()

	maxExp := model.UnderstandExp + model.MasteredExp

	if exp < maxExp {
		numOfAnswer := len(input.Answers)

		englishItem.SetExp(exp + uint(numOfAnswer))
		englishItem.CheckExp()

		u.englishItemRepo.Update(englishItem)
	}

	outputs := []*model.Output{}

	for _, answer := range answers {
		output := model.NewOutput("", englishItemId, content, answer.Question, answer.Answer, "", answer.Index, time.Time{})
		outputs = append(outputs, output)
	}

	if err := u.chatAIAPI.GetAdvice(outputs); err != nil {
		return nil, err
	}

	adviceList := []*dto.Advice{}

	for _, output := range outputs {
		advice := &dto.Advice{
			Index:  output.Index(),
			Advice: output.Advice(),
		}

		adviceList = append(adviceList, advice)
	}

	answerQuestionsOutput := &dto.AnswerQuestionsOutput{
		AdviceList: adviceList,
	}

	return answerQuestionsOutput, nil
}
