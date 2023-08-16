package userusecase

import (
	"english/config"
	"english/myerror"
	"english/src/domain/usermodel"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginUsecase interface {
	Login(req *UserLoginRequest) (string, error)
}

type UserJWTLoginUsecase struct {
	ur usermodel.UserRepository
}

func NewUserJWTLoginUsecase(ur usermodel.UserRepository) UserLoginUsecase {
	return &UserJWTLoginUsecase{
		ur: ur,
	}
}

func (uu *UserJWTLoginUsecase) Login(req *UserLoginRequest) (string, error) {
	user, err := uu.ur.FindByEmail(req.Email)
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

	// JWTの作成
	// 有効期限1日
	jwtToken, err := createJWT(user.Id(), 60*60*24)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func createJWT(userId int, expSec int) (string, error) {
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

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
