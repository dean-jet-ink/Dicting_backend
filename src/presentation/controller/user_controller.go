package controller

import (
	"english/config"
	"english/myerror"
	"english/src/usecase/userusecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

type UserGinController struct {
	ulu userusecase.UserLoginUsecase
}

func NewUserGinController(ulu userusecase.UserLoginUsecase) UserController {
	return &UserGinController{
		ulu: ulu,
	}
}

func (uc *UserGinController) Login(c *gin.Context) {
	req := &userusecase.UserLoginRequest{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := uc.ulu.Login(req)
	if err != nil {
		if errors.Is(myerror.ErrRecordNotFound, err) || errors.Is(myerror.ErrMismatchedPassword, err) {
			c.JSON(http.StatusUnauthorized, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	// JWTをcookieへセット
	// 有効期限 1日
	// 開発環境では、postman等での開発用にSecureをfalseに指定
	c.SetCookie("token", jwtToken, 60*60*24, "/", config.APIDomain(), config.GoEnv() != "dev", true)
	// クロスオリジン許可
	c.SetSameSite(http.SameSiteNoneMode)

	c.JSON(http.StatusOK, jwtToken)
}

func (uc *UserGinController) Logout(c *gin.Context) {
	// cookieのJWTトークン削除
	// MaxAgeに負の数を指定 = 即時削除
	c.SetCookie("token", "", -1, "/", config.APIDomain(), config.GoEnv() != "dev", true)
}
