package api

import "english/cmd/domain/model"

type ChatAIAPI interface {
	GetTranslations(englishItem *model.EnglishItem) error
	GetExamples(englishItem *model.EnglishItem) error
	GetTranslation(content string) (string, error)
	GetExplanation(content string) (string, error)
	GetExample(content string) (*model.Example, error)
	GetQuestion(content string) (string, error)
	GetAdvice(answers []*model.Output) error
}
