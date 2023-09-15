package usecase

import (
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
)

type GetUserUsecase interface {
	GetUser(id string) (*dto.UserResponse, error)
}

type GetUserUsecaseImpl struct {
	ur repository.UserRepository
}

func NewGetUserUsecase(ur repository.UserRepository) GetUserUsecase {
	return &GetUserUsecaseImpl{
		ur: ur,
	}
}

func (u *GetUserUsecaseImpl) GetUser(id string) (*dto.UserResponse, error) {
	user, err := u.ur.FindById(id)
	if err != nil {
		return nil, err
	}

	res := &dto.UserResponse{
		Email: user.Email(),
		Name:  user.Name(),
		Image: user.ProfileImageURL(),
	}

	return res, nil
}
