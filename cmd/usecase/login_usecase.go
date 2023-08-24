package usecase

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
	"english/myerror"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type LoginUsecase interface {
	Login(req *dto.LoginRequest, isSSO bool) (string, error)
}

type LoginUsecaseImpl struct {
	ur repository.UserRepository
}

func NewLoginUsecase(ur repository.UserRepository) LoginUsecase {
	return &LoginUsecaseImpl{
		ur: ur,
	}
}

func (uu *LoginUsecaseImpl) Login(req *dto.LoginRequest, isSSO bool) (string, error) {
	user := &model.User{}
	var err error

	if isSSO {
		user, err = uu.ur.FindByIssAndSub(req.Iss, req.Sub)
		if err != nil {
			return "", err
		}
	} else {
		user, err = uu.ur.FindByEmail(req.Email)
		if err != nil {
			return "", err
		}

		// ハッシュ化して保存しているパスワードと比較
		err = bcrypt.CompareHashAndPassword([]byte(user.Password()), []byte(req.Password))
		if err != nil {
			if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) {
				return "", myerror.ErrMismatchedPassword
			}

			return "", err
		}
	}

	// JWTの作成
	// 有効期限1日
	jwtToken, err := user.CreateJWT(60 * 60 * 24)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
