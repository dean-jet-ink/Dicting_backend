package usecase

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
	"english/cmd/usecase/validator"
	"english/lib"
)

type CreateEnglishItemUsecase interface {
	Create(req *dto.CreateEnglishItemRequest) error
}

type CreateEnglishItemUsecaseImpl struct {
	englishItemRepo repository.EnglishItemRepository
	fileStorageRepo repository.FileStorageRepository
	v               validator.EnglishItemValidator
}

func NewCreateEnglishItemUsecase(englishItemRepo repository.EnglishItemRepository, fileStorageRepo repository.FileStorageRepository, v validator.EnglishItemValidator) CreateEnglishItemUsecase {
	return &CreateEnglishItemUsecaseImpl{
		englishItemRepo: englishItemRepo,
		fileStorageRepo: fileStorageRepo,
		v:               v,
	}
}

func (u *CreateEnglishItemUsecaseImpl) Create(req *dto.CreateEnglishItemRequest) error {
	imgs := []*model.Img{}
	for _, reqImg := range req.Imgs {
		img := model.NewImg(reqImg.Id, reqImg.URL, reqImg.IsThumbnail)

		ulid, err := lib.GenerateULID()
		if err != nil {
			return err
		}
		img.SetId(ulid)

		imgs = append(imgs, img)
	}

	if err := u.fileStorageRepo.UploadImgs(imgs, nil); err != nil {
		return err
	}

	examples := []*model.Example{}
	for _, example := range req.Examples {
		ulid, err := lib.GenerateULID()
		if err != nil {
			return err
		}

		exampleDomain := model.NewExample(ulid, example.Example, example.Translation)
		examples = append(examples, exampleDomain)
	}

	ulid, err := lib.GenerateULID()
	if err != nil {
		return err
	}

	englishItem := model.NewEnglishItem(ulid, req.Content, req.Translations, req.EnExplanation, examples, imgs, req.UserId, model.Learning, 0)

	if err := u.v.EnglishItemValidate(englishItem); err != nil {
		return err
	}

	if err = u.englishItemRepo.Create(englishItem); err != nil {
		return err
	}

	return nil
}
