package controller

import (
	"english/cmd/usecase"
	"english/cmd/usecase/dto"
	"english/config"
	"english/myerror"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

type UserGinController struct {
	su usecase.SignupUsecase
	lu usecase.LoginUsecase
}

func NewUserGinController(su usecase.SignupUsecase, lu usecase.LoginUsecase) UserController {
	return &UserGinController{
		su: su,
		lu: lu,
	}
}

func (uc *UserGinController) Signup(c *gin.Context) {
	req := &dto.SignupRequest{}
	if err := c.BindJSON(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := uc.su.Signup(req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	uc.setJWT(c, jwtToken)

	c.JSON(http.StatusOK, jwtToken)
}

func (uc *UserGinController) Login(c *gin.Context) {
	req := &dto.LoginRequest{}
	if err := c.BindJSON(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := uc.lu.Login(req)
	if err != nil {
		if errors.Is(myerror.ErrRecordNotFound, err) || errors.Is(myerror.ErrMismatchedPassword, err) {
			c.JSON(http.StatusUnauthorized, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	uc.setJWT(c, jwtToken)

	c.JSON(http.StatusOK, jwtToken)
}

func (uc *UserGinController) Logout(c *gin.Context) {
	// cookieのJWTトークン削除
	// MaxAgeに負の数を指定 = 即時削除
	c.SetCookie("token", "", -1, "/", config.APIDomain(), config.GoEnv() != "dev", true)
}

func (uc *UserGinController) setJWT(c *gin.Context, jwtToken string) {
	// JWTをcookieへセット
	// 有効期限 1日
	// 開発環境では、postman等での開発用にSecureをfalseに指定
	c.SetCookie("token", jwtToken, 60*60*24, "/", config.APIDomain(), config.GoEnv() != "dev", true)
	// クロスオリジン許可
	c.SetSameSite(http.SameSiteNoneMode)
}
