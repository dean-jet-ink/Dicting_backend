package usecase

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
	"english/lib"

	"golang.org/x/crypto/bcrypt"
)

type SignupUsecase interface {
	Signup(req *dto.SignupRequest, isSSO bool) (string, error)
}

type SignupUsecaseImpl struct {
	ur repository.UserRepository
}

func NewSignupUsecase(ur repository.UserRepository) SignupUsecase {
	return &SignupUsecaseImpl{
		ur: ur,
	}
}

func (su *SignupUsecaseImpl) Signup(req *dto.SignupRequest, isSSO bool) (string, error) {
	ulid, err := lib.GenerateULID()
	if err != nil {
		return "", err
	}

	user := model.NewUser(ulid, req.Email, "", req.Name, "")
	if !isSSO {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			return "", err
		}
		pass := string(hash)
		user.SetPassword(pass)
	} else {
		user.SetIss(req.Iss)
		user.SetSub(req.Sub)
	}

	if err := su.ur.Create(user); err != nil {
		return "", err
	}

	jwtToken, err := user.CreateJWT(60 * 60 * 24)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
