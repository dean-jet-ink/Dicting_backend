package repository

import "english/cmd/domain/model"

type EnglishItemRepository interface {
	Create(englishItem *model.EnglishItem) error
	FindByUserIdAndContent(userId, content string) ([]*model.EnglishItem, error)
}
