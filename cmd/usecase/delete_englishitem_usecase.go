package usecase

import (
	"english/cmd/domain/repository"
)

type DeleteEnglishItemUsecase interface {
	Delete(id string) error
}

type DeleteEnglishItemUsecaseImpl struct {
	englishItemRepo repository.EnglishItemRepository
	fileStorageRepo repository.FileStorageRepository
}

func NewDeleteEnglishItemUsecase(englishItemRepo repository.EnglishItemRepository, fileStorageRepo repository.FileStorageRepository) DeleteEnglishItemUsecase {
	return &DeleteEnglishItemUsecaseImpl{
		englishItemRepo: englishItemRepo,
		fileStorageRepo: fileStorageRepo,
	}
}

func (u *DeleteEnglishItemUsecaseImpl) Delete(id string) error {
	imgs, err := u.englishItemRepo.FindImgsByEnglishItemId(id)
	if err != nil {
		return err
	}

	if err := u.fileStorageRepo.DeleteImgs(imgs); err != nil {
		return err
	}

	if err := u.englishItemRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
