package usecase

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
)

type GetEnglishItemUsecase interface {
	GetByUserIdAndContent(userId, content string) (*dto.GetEnglishItemResponse, error)
}

type GetEnglishItemUsecaseImpl struct {
	er repository.EnglishItemRepository
}

func NewGetEnglishItemUsecase(er repository.EnglishItemRepository) GetEnglishItemUsecase {
	return &GetEnglishItemUsecaseImpl{
		er: er,
	}
}

func (u *GetEnglishItemUsecaseImpl) GetByUserIdAndContent(userId, content string) (*dto.GetEnglishItemResponse, error) {
	englishItems, err := u.er.FindByUserIdAndContent(userId, content)
	if err != nil {
		return nil, err
	}

	englishItemDTOs := []*dto.EnglishItem{}
	for _, e := range englishItems {
		examples := u.exampleDomainsToDTOs(e.Examples())
		imgs := u.imgDomainsToDTOs(e.Imgs())
		englishItemDTO := &dto.EnglishItem{
			Id:            e.Id(),
			Content:       e.Content(),
			Translations:  e.Translations(),
			EnExplanation: e.EnExplanation(),
			Examples:      examples,
			Imgs:          imgs,
		}

		englishItemDTOs = append(englishItemDTOs, englishItemDTO)
	}

	resp := &dto.GetEnglishItemResponse{
		EnglishItems: englishItemDTOs,
	}

	return resp, nil
}

func (u *GetEnglishItemUsecaseImpl) exampleDomainsToDTOs(examples []*model.Example) []*dto.Example {
	exampleDTOs := []*dto.Example{}
	for _, example := range examples {
		exampleDTO := &dto.Example{
			Id:          example.Id,
			Example:     example.Example,
			Translation: example.Translation,
		}
		exampleDTOs = append(exampleDTOs, exampleDTO)
	}

	return exampleDTOs
}

func (u *GetEnglishItemUsecaseImpl) imgDomainsToDTOs(imgs []*model.Img) []*dto.Img {
	imgDTOs := []*dto.Img{}
	for _, img := range imgs {
		imgDTO := &dto.Img{
			Id:  img.Id(),
			URL: img.URL(),
		}
		imgDTOs = append(imgDTOs, imgDTO)
	}

	return imgDTOs
}
