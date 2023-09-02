package usecase

import (
	"english/algo"
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
)

type CreateEnglishItemUsecase interface {
	Create(req *dto.CreateEnglishItemRequest) (*dto.CreateEnglishItemResponse, error)
}

type CreateEnglishItemUsecaseImpl struct {
	englishItemRepo repository.EnglishItemRepository
	fileStorageRepo repository.FileStorageRepository
}

func NewCreateEnglishItemUsecase(englishItemRepo repository.EnglishItemRepository, fileStorageRepo repository.FileStorageRepository) CreateEnglishItemUsecase {
	return &CreateEnglishItemUsecaseImpl{
		englishItemRepo: englishItemRepo,
		fileStorageRepo: fileStorageRepo,
	}
}

func (u *CreateEnglishItemUsecaseImpl) Create(req *dto.CreateEnglishItemRequest) (*dto.CreateEnglishItemResponse, error) {
	newURLs, err := u.fileStorageRepo.UploadFromURLs(req.ImgURLs, nil)
	if err != nil {
		return nil, err
	}

	ulid, err := algo.GenerateULID()
	if err != nil {
		return nil, err
	}

	examples := []*model.Example{}
	for _, example := range req.Examples {
		ulid, err := algo.GenerateULID()
		if err != nil {
			return nil, err
		}

		exampleDomain := model.NewExample(ulid, example.Example, example.Translation)
		examples = append(examples, exampleDomain)
	}

	imgs := []*model.Img{}
	for _, url := range newURLs {
		ulid, err := algo.GenerateULID()
		if err != nil {
			return nil, err
		}

		imgDomain := model.NewImg(ulid, url)
		imgs = append(imgs, imgDomain)
	}

	englishItem := model.NewEnglishItem(ulid, req.Content, req.JaTranslations, req.EnExplanation, examples, imgs, req.UserId)

	if err = u.englishItemRepo.Create(englishItem); err != nil {
		return nil, err
	}

	exampleDTOs := []*dto.Example{}
	for _, example := range englishItem.Examples() {
		exampleDTO := *&dto.Example{
			Id:          example.Id,
			Example:     example.Example,
			Translation: example.Translation,
		}
		exampleDTOs = append(exampleDTOs, &exampleDTO)
	}

	imgDTOs := []*dto.Img{}
	for _, img := range englishItem.Imgs() {
		imgDTO := &dto.Img{
			Id:  img.Id(),
			URL: img.URL(),
		}
		imgDTOs = append(imgDTOs, imgDTO)
	}

	resp := &dto.CreateEnglishItemResponse{
		Id:             englishItem.Id(),
		Content:        englishItem.Content(),
		JaTranslations: englishItem.Translations(),
		EnExplanation:  englishItem.EnExplanation(),
		Imgs:           imgDTOs,
		Examples:       exampleDTOs,
	}

	return resp, nil
}
