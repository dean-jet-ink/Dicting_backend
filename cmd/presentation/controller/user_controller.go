package controller

import (
	"english/cmd/presentation/errhandle"
	"english/cmd/usecase"
	"english/cmd/usecase/dto"
	"english/config"
	"english/myerror"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RedirectOAuthConsent(c *gin.Context)
	OAuthCallback(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateProfile(c *gin.Context)
	UpdateProfileImg(c *gin.Context)
}

type UserGinController struct {
	su  usecase.SignupUsecase
	lu  usecase.LoginUsecase
	ssu usecase.SSOAuthUsecase
	gu  usecase.GetUserUsecase
	uu  usecase.UpdateUserUsecase
	upu usecase.UpdateProfileImgUsecase
}

func NewUserGinController(su usecase.SignupUsecase, lu usecase.LoginUsecase, ssu usecase.SSOAuthUsecase, gu usecase.GetUserUsecase, uu usecase.UpdateUserUsecase, upu usecase.UpdateProfileImgUsecase) UserController {
	return &UserGinController{
		su:  su,
		lu:  lu,
		ssu: ssu,
		gu:  gu,
		uu:  uu,
		upu: upu,
	}
}

func (uc *UserGinController) Signup(c *gin.Context) {
	req := &dto.SignupRequest{}
	if err := c.BindJSON(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrBindingFailure, c)
		return
	}

	jwtToken, err := uc.su.Signup(req, false)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	uc.setJWT(c, jwtToken)

	c.JSON(http.StatusOK, jwtToken)
}

func (uc *UserGinController) Login(c *gin.Context) {
	req := &dto.LoginRequest{}
	if err := c.BindJSON(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrBindingFailure, c)
		return
	}

	if err := Validate(req); err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	jwtToken, err := uc.lu.Login(req, false)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	uc.setJWT(c, jwtToken)

	c.JSON(http.StatusOK, jwtToken)
}

func (uc *UserGinController) Logout(c *gin.Context) {
	// cookieのJWTトークン削除
	// MaxAgeに負の数を指定 = 即時削除
	uc.deleteCookie(c, "token", "/")
}

func (uc *UserGinController) RedirectOAuthConsent(c *gin.Context) {
	req := &dto.RedirectOAuthConsentRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrBindingFailure, c)
		return
	}

	if err := Validate(req); err != nil {
		errhandle.HandleValidationError(err, c)
		return
	}

	resp, err := uc.ssu.RedirectOAuthConsent(req)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	maxAge := 60 * 60 * 24
	path := "/auth/"
	domain := config.APIDomain()
	isSecure := config.GoEnv() != "dev"

	// callbackで使用する検証用state
	c.SetCookie("oauth_state", resp.State, maxAge, path, domain, isSecure, true)
	// callbackで使用するidPName
	c.SetCookie("idp_name", req.IdPName, maxAge, path, domain, isSecure, true)
	// ログインフラグ
	c.SetCookie("is_login", strconv.FormatBool(req.IsLogin), maxAge, path, domain, isSecure, true)
	c.SetSameSite(http.SameSiteNoneMode)

	c.Redirect(http.StatusFound, resp.RedirectURL)
}

func (uc *UserGinController) OAuthCallback(c *gin.Context) {
	errRedirectURL := fmt.Sprintf("%v/login", config.FrontEndURL())

	req := &dto.CallbackRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		errhandle.HandleErrorRedirect(myerror.ErrBindingFailure, c, errRedirectURL)
		return
	}

	state, err := c.Cookie("oauth_state")
	if err != nil {
		errhandle.HandleErrorRedirect(fmt.Errorf("%v: %w", myerror.ErrCookieExpired, errors.New("oauth_state")), c, errRedirectURL)
		return
	}
	idPName, err := c.Cookie("idp_name")
	if err != nil {
		errhandle.HandleErrorRedirect(fmt.Errorf("%v: %w", myerror.ErrCookieExpired, errors.New("idp_name")), c, errRedirectURL)
		return
	}
	isLoginStr, err := c.Cookie("is_login")
	if err != nil {
		errhandle.HandleErrorRedirect(fmt.Errorf("%v: %w", myerror.ErrCookieExpired, errors.New("is_login")), c, errRedirectURL)
		return
	}

	isLogin, err := strconv.ParseBool(isLoginStr)
	if err != nil {
		errhandle.HandleErrorRedirect(fmt.Errorf("%v: %w", myerror.ErrParsingFailure, err), c, errRedirectURL)
		return
	}

	req.CookieState = state
	req.IdpName = idPName

	if err := Validate(req); err != nil {
		errhandle.HandleValidationError(err, c)
		return
	}

	resp, err := uc.ssu.Callback(req)
	if err != nil {
		errhandle.HandleErrorRedirect(err, c, errRedirectURL)
		return
	}

	var jwtToken string
	if isLogin {
		req := &dto.LoginRequest{
			Email:    resp.Email,
			Password: "",
			Iss:      resp.Iss,
			Sub:      resp.Sub,
		}

		jwtToken, err = uc.lu.Login(req, true)
		if err != nil {
			errhandle.HandleErrorRedirect(err, c, errRedirectURL)
			return
		}
	} else {
		if !resp.Verified {
			errhandle.HandleErrorRedirect(myerror.ErrUnverifiedEmail, c, errRedirectURL)
			return
		}

		req := &dto.SignupRequest{
			Email: resp.Email,
			Name:  resp.Name,
			Iss:   resp.Iss,
			Sub:   resp.Sub,
		}

		jwtToken, err = uc.su.Signup(req, true)
		if err != nil {
			// 指定のIdPでサインアップ済みの場合
			if errors.Is(myerror.ErrDuplicatedKey, err) {

				req := &dto.LoginRequest{
					Email:    resp.Email,
					Password: "",
					Iss:      resp.Iss,
					Sub:      resp.Sub,
				}

				jwtToken, err = uc.lu.Login(req, true)
				if err != nil {
					errhandle.HandleErrorRedirect(err, c, errRedirectURL)
					return
				}
			} else {
				errhandle.HandleErrorRedirect(err, c, errRedirectURL)
				return
			}
		}
	}

	uc.setJWT(c, jwtToken)

	// 不要なoauth_stateとidp_nameを削除
	uc.deleteCookie(c, "oauth_state", "/auth/")
	uc.deleteCookie(c, "idp_name", "/auth/")
	uc.deleteCookie(c, "is_login", "/auth/")

	redirectURL := fmt.Sprintf("%v/", config.FrontEndURL())
	c.Redirect(http.StatusFound, redirectURL)
}

func (uc *UserGinController) setJWT(c *gin.Context, jwtToken string) {
	// JWTをcookieへセット
	// 有効期限 1日
	// 開発環境では、postman等での開発用にSecureをfalseに指定
	c.SetCookie("token", jwtToken, 60*60*24, "/", config.APIDomain(), config.GoEnv() != "dev", true)
	// クロスオリジン許可
	c.SetSameSite(http.SameSiteNoneMode)
}

func (uc *UserGinController) deleteCookie(c *gin.Context, name, path string) {
	c.SetCookie(name, "", -1, path, config.APIDomain(), config.GoEnv() != "dev", true)
}

func (uc *UserGinController) GetUser(c *gin.Context) {

	id, err := userId(c)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	res, err := uc.gu.GetUser(id)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (uc *UserGinController) UpdateProfile(c *gin.Context) {
	req := &dto.UpdateUserRequest{}
	if err := c.Bind(req); err != nil {
		errhandle.HandleErrorJSON(myerror.ErrBindingFailure, c)
		return
	}

	id, err := userId(c)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}
	req.Id = id

	if err := Validate(req); err != nil {
		errhandle.HandleValidationError(err, c)
		return
	}

	resp, err := uc.uu.Update(req)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (uc *UserGinController) UpdateProfileImg(c *gin.Context) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		errhandle.HandleErrorJSON(myerror.ErrFormFileNotFound, c)
		return
	}

	id, err := userId(c)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	input := &dto.UpdateProfileImgInput{
		Id:         id,
		FileHeader: fileHeader,
	}

	output, err := uc.upu.Update(input)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	updateUserReq := &dto.UpdateUserRequest{
		Id:    id,
		Image: output.ProfileImgURL,
	}

	updateUserResp, err := uc.uu.Update(updateUserReq)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		return
	}

	c.JSON(http.StatusOK, updateUserResp)
}
