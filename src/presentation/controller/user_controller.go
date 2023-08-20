package controller

import (
	"english/config"
	"english/myerror"
	"english/src/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RedirectOAuthConsent(c *gin.Context)
	OAuthCallback(c *gin.Context)
}

type UserGinController struct {
	lu usecase.LoginUsecase
	su usecase.SSOLoginUsecase
}

func NewUserGinController(lu usecase.LoginUsecase, su usecase.SSOLoginUsecase) UserController {
	return &UserGinController{
		lu: lu,
		su: su,
	}
}

func (uc *UserGinController) Login(c *gin.Context) {
	req := &usecase.LoginRequest{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := uc.lu.Login(req, false)
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

func (uc *UserGinController) RedirectOAuthConsent(c *gin.Context) {
	req := &usecase.RedirectOAuthConsentRequest{}
	if err := c.BindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	redirectOut, err := uc.su.RedirectOAuthConsent(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	maxAge := 60 * 60 * 24
	path := "/auth/"
	domain := config.APIDomain()
	isSecure := config.GoEnv() != "dev"

	// callbackで使用する検証用state
	c.SetCookie("oauth_state", redirectOut.State(), maxAge, path, domain, isSecure, true)

	// callbackで使用するidPName
	c.SetCookie("idp_name", req.IdPName, maxAge, path, domain, isSecure, true)

	c.Redirect(http.StatusFound, redirectOut.RedirectURL())
}

func (uc *UserGinController) OAuthCallback(c *gin.Context) {
	req := &usecase.CallbackRequest{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	state, stateErr := c.Cookie("oauth_state")
	idPName, idPErr := c.Cookie("idp_name")
	if stateErr != nil || idPErr != nil {
		c.JSON(http.StatusBadRequest, "cookie expired")
		return
	}

	req.CookieState = state
	req.IdpName = idPName

	if err := Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	callbackOut, err := uc.su.Callback(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !callbackOut.Verified {
		c.JSON(http.StatusForbidden, "Email is unverified")
		return
	}

	jwtToken, err := uc.lu.Login(&usecase.LoginRequest{
		Email:    callbackOut.Email,
		Password: "",
	}, true)

	uc.setJWT(c, jwtToken)

	c.JSON(http.StatusOK, jwtToken)
}

func (uc *UserGinController) setJWT(c *gin.Context, jwtToken string) {
	// JWTをcookieへセット
	// 有効期限 1日
	// 開発環境では、postman等での開発用にSecureをfalseに指定
	c.SetCookie("token", jwtToken, 60*60*24, "/", config.APIDomain(), config.GoEnv() != "dev", true)
	// クロスオリジン許可
	c.SetSameSite(http.SameSiteNoneMode)
}
