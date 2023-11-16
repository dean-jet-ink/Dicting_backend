package usecase

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
	"english/cmd/usecase/validator"
	"english/lib"
)

type UpdateEnglishItemUsecase interface {
	Update(req *dto.UpdateEnglishItemRequest) error
}

type UpdateEnglishItemUsecaseImpl struct {
	englishItemRepo repository.EnglishItemRepository
	fileStorageRepo repository.FileStorageRepository
	v               validator.EnglishItemValidator
}

func NewUpdateEnglishItemUsecase(englishItemRepo repository.EnglishItemRepository, fileStorageRepo repository.FileStorageRepository, v validator.EnglishItemValidator) UpdateEnglishItemUsecase {
	return &UpdateEnglishItemUsecaseImpl{
		englishItemRepo: englishItemRepo,
		fileStorageRepo: fileStorageRepo,
		v:               v,
	}
}

func (u *UpdateEnglishItemUsecaseImpl) Update(req *dto.UpdateEnglishItemRequest) error {
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

	englishItemId := req.Id

	preImgs, err := u.englishItemRepo.FindImgsByEnglishItemId(englishItemId)
	if err != nil {
		return err
	}

	if err := u.fileStorageRepo.UploadImgs(imgs, preImgs); err != nil {
		return err
	}

	if err = u.englishItemRepo.DeleteImgsByEnglishItemId(englishItemId); err != nil {
		return err
	}

	if err = u.englishItemRepo.DeleteExampleByEnglishItemId(englishItemId); err != nil {
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

	englishItem := model.NewEnglishItem(req.Id, req.Content, req.Translations, req.EnExplanation, examples, imgs, req.UserId, model.Proficiency(req.Proficiency), req.Exp)

	if err := u.v.EnglishItemValidate(englishItem); err != nil {
		return err
	}

	if err := u.englishItemRepo.Update(englishItem); err != nil {
		return err
	}

	return nil
}
