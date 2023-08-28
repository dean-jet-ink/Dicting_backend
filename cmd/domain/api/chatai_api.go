package api

import "english/cmd/domain/model"

type ChatAIAPI interface {
	GetTranslation(englishItem *model.EnglishItem) error
	GetExample(englishItem *model.EnglishItem) error
}
