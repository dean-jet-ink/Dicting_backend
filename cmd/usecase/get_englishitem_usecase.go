package usecase

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
)

type GetEnglishItemUsecase interface {
	GetEnglishItemInfoByUserId(userId string) (*dto.GetEnglishItemsResponse, error)
	GetById(id string) (*dto.GetEnglishItemResponse, error)
}

type GetEnglishItemUsecaseImpl struct {
	er repository.EnglishItemRepository
}

func NewGetEnglishItemUsecase(er repository.EnglishItemRepository) GetEnglishItemUsecase {
	return &GetEnglishItemUsecaseImpl{
		er: er,
	}
}

func (u *GetEnglishItemUsecaseImpl) GetEnglishItemInfoByUserId(userId string) (*dto.GetEnglishItemsResponse, error) {
	englishItems, err := u.er.FindEnglishItemInfosByUserId(userId)
	if err != nil {
		return nil, err
	}

	resp := &dto.GetEnglishItemsResponse{
		EnglishItems: []*dto.EnglishItem{},
	}

	for _, item := range englishItems {
		if len(item.Imgs()) == 0 {
			item.SetImgs([]*model.Img{model.NewImg("", "", false)})
		}

		dto := &dto.EnglishItem{
			Id:            item.Id(),
			Content:       item.Content(),
			Translations:  item.Translations(),
			EnExplanation: item.EnExplanation(),
			Img:           item.Imgs()[0].URL(),
			Proficiency:   item.Proficiency(),
			Exp:           item.Exp(),
		}

		resp.EnglishItems = append(resp.EnglishItems, dto)
	}

	return resp, nil
}

func (u *GetEnglishItemUsecaseImpl) GetById(id string) (*dto.GetEnglishItemResponse, error) {
	englishItem, err := u.er.FindById(id)
	if err != nil {
		return nil, err
	}

	resp := &dto.GetEnglishItemResponse{
		Id:            englishItem.Id(),
		Content:       englishItem.Content(),
		Translations:  englishItem.Translations(),
		EnExplanation: englishItem.EnExplanation(),
		Examples:      u.exampleDomainsToDTOs(englishItem.Examples()),
		Imgs:          u.imgDomainsToDTOs(englishItem.Imgs()),
		Proficiency:   englishItem.Proficiency(),
		Exp:           englishItem.Exp(),
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
			Id:          img.Id(),
			URL:         img.URL(),
			IsThumbnail: img.IsThumbnail(),
		}
		imgDTOs = append(imgDTOs, imgDTO)
	}

	return imgDTOs
}
