package usecase

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
)

type UpdateUserUsecase interface {
	Update(req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
}

type UpdateUserProfileUsecase struct {
	ur repository.UserRepository
}

func NewUpdateUserProfileUsecase(ur repository.UserRepository) UpdateUserUsecase {
	return &UpdateUserProfileUsecase{
		ur: ur,
	}
}

func (uu *UpdateUserProfileUsecase) Update(req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	user := model.NewUser(req.Id, req.Email, "", req.Name, req.Image)

	if err := uu.ur.Update(user); err != nil {
		return nil, err
	}

	resp := &dto.UpdateUserResponse{
		Email: user.Email(),
		Name:  user.Name(),
		Image: user.ProfileImageURL(),
	}

	return resp, nil
}
