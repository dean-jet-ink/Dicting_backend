package usecase

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
	"english/lib"
)

type UpdateProfileImgUsecase interface {
	Update(req *dto.UpdateProfileImgInput) (*dto.UpdateProfileImgOutput, error)
}

type UpdateProfileImgUsecaseImpl struct {
	ur repository.UserRepository
	fr repository.FileStorageRepository
}

func NewUpdateProfileImgUsecase(ur repository.UserRepository, fr repository.FileStorageRepository) UpdateProfileImgUsecase {
	return &UpdateProfileImgUsecaseImpl{
		ur: ur,
		fr: fr,
	}
}

func (u *UpdateProfileImgUsecaseImpl) Update(input *dto.UpdateProfileImgInput) (*dto.UpdateProfileImgOutput, error) {

	user, err := u.ur.FindById(input.Id)
	if err != nil {
		return nil, err
	}

	preFile := user.ProfileImageURL()

	file, err := input.FileHeader.Open()
	if err != nil {
		return nil, err
	}

	ulid, err := lib.GenerateULID()
	if err != nil {
		return nil, err
	}

	imgFile := &model.ImgFile{
		Body:     file,
		FileName: ulid,
	}

	if err := u.fr.Upload(imgFile, preFile); err != nil {
		return nil, err
	}

	resp := &dto.UpdateProfileImgOutput{
		ProfileImgURL: imgFile.URL,
	}

	return resp, nil
}
