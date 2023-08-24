package usecase

import (
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
	"english/config"
)

type UpdateProfileImgUsecase interface {
	Update(req *dto.UpdateProfileImgRequest) (*dto.UpdateProfileImgResponse, error)
}

type UpdateProfileImgUsecaseImpl struct {
	ur repository.UserRepository
}

func NewUpdateProfileImgUsecase(ur repository.UserRepository) UpdateProfileImgUsecase {
	return &UpdateProfileImgUsecaseImpl{
		ur: ur,
	}
}

func (upu *UpdateProfileImgUsecaseImpl) Update(req *dto.UpdateProfileImgRequest) (*dto.UpdateProfileImgResponse, error) {
	resp := &dto.UpdateProfileImgResponse{}

	user, err := upu.ur.FindById(req.Id)
	if err != nil {
		return nil, err
	}

	if config.GoEnv() == "dev" {
		resp.ProfileImgURL = user.ProfileImageURL()
		return resp, nil
	}

	// S3またはCloud Storageでの処理

	return resp, nil
}
