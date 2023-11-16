package repository

import "english/cmd/domain/model"

type EnglishItemRepository interface {
	Create(englishItem *model.EnglishItem) error
	Update(englishItem *model.EnglishItem) error
	Delete(englishItemId string) error
	FindEnglishItemInfosByUserId(userId string) ([]*model.EnglishItem, error)
	FindById(id string) (*model.EnglishItem, error)
	FindImgsByEnglishItemId(englishItemId string) ([]*model.Img, error)
	DeleteImgsByEnglishItemId(englishItemId string) error
	DeleteExampleByEnglishItemId(englishItemId string) error
}
