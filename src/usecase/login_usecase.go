package usecase

import (
	"english/config"
	"english/myerror"
	"english/src/domain/usermodel"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginUsecase interface {
	Login(req *LoginRequest, isSSO bool) (string, error)
}

type JWTLoginUsecase struct {
	ur usermodel.UserRepository
}

func NewJWTLoginUsecase(ur usermodel.UserRepository) LoginUsecase {
	return &JWTLoginUsecase{
		ur: ur,
	}
}

func (lu *JWTLoginUsecase) Login(req *LoginRequest, isSSO bool) (string, error) {
	user, err := lu.ur.FindByEmail(req.Email)
	if err != nil {
		return "", err
	}

	if !isSSO {
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
	jwtToken, err := createJWT(user.Id(), 60*60*24)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func createJWT(userId string, expSec int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Second * time.Duration(expSec)).Unix(),
	})
	// 環境変数のSECRETを使用し署名
	jwtToken, err := token.SignedString([]byte(config.Secret()))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gte=4,lt=30"`
}
