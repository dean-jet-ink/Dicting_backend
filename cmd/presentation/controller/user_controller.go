package controller

import (
	"english/algo"
	"english/cmd/usecase"
	"english/cmd/usecase/dto"
	"english/config"
	"english/myerror"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RedirectOAuthConsent(c *gin.Context)
	OAuthCallback(c *gin.Context)
	UpdateProfile(c *gin.Context)
	UpdateProfileImg(c *gin.Context)
}

type UserGinController struct {
	su  usecase.SignupUsecase
	lu  usecase.LoginUsecase
	ssu usecase.SSOAuthUsecase
	uu  usecase.UpdateUserUsecase
	upu usecase.UpdateProfileImgUsecase
}

func NewUserGinController(su usecase.SignupUsecase, lu usecase.LoginUsecase, ssu usecase.SSOAuthUsecase, uu usecase.UpdateUserUsecase, upu usecase.UpdateProfileImgUsecase) UserController {
	return &UserGinController{
		su:  su,
		lu:  lu,
		ssu: ssu,
		uu:  uu,
		upu: upu,
	}
}

func (uc *UserGinController) Signup(c *gin.Context) {
	req := &dto.SignupRequest{}
	if err := c.BindJSON(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := uc.su.Signup(req, false)
	if err != nil {
		if errors.Is(myerror.ErrDuplicatedKey, err) {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
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
		log.Println(err)
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

		log.Println(err)
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
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := uc.ssu.RedirectOAuthConsent(req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
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
	req := &dto.CallbackRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	state, err1 := c.Cookie("oauth_state")
	idPName, err2 := c.Cookie("idp_name")
	isLoginStr, err3 := c.Cookie("is_login")
	if err1 != nil || err2 != nil || err3 != nil {
		log.Println(errors.New("cookie expired"))
		c.JSON(http.StatusBadRequest, "cookie expired")
		return
	}
	isLogin, err := strconv.ParseBool(isLoginStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	req.CookieState = state
	req.IdpName = idPName

	if err := Validate(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := uc.ssu.Callback(req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
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
			if errors.Is(myerror.ErrRecordNotFound, err) {
				message := "user is not registered through SSO"
				log.Println(errors.New(message))
				c.JSON(http.StatusUnauthorized, message)
			} else {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, err.Error())
			}
			return
		}
	} else {
		if !resp.Verified {
			message := "email is unverified"
			log.Println(errors.New(message))
			c.JSON(http.StatusForbidden, message)
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
			if errors.Is(myerror.ErrDuplicatedKey, err) {
				log.Println(err)
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			log.Println(err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	uc.setJWT(c, jwtToken)

	// 不要なoauth_stateとidp_nameを削除
	uc.deleteCookie(c, "oauth_state", "/auth/")
	uc.deleteCookie(c, "idp_name", "/auth/")
	uc.deleteCookie(c, "is_login", "/auth/")

	if isLogin {
		c.String(http.StatusOK, jwtToken)
	} else {
		c.String(http.StatusCreated, jwtToken)
	}
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

func (uc *UserGinController) UpdateProfile(c *gin.Context) {
	req := &dto.UpdateUserRequest{}
	if err := c.Bind(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err := userId(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	req.Id = id

	if err := Validate(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := uc.uu.Update(req)
	if err != nil {
		if errors.Is(myerror.ErrDuplicatedKey, err) {
			c.JSON(http.StatusBadRequest, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (uc *UserGinController) UpdateProfileImg(c *gin.Context) {
	file, err := c.FormFile("profile_img")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err := userId(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	updateProfileImgReq := &dto.UpdateProfileImgRequest{
		Id:   id,
		File: file,
	}
	updateProfileImgResp, err := uc.upu.Update(updateProfileImgReq)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if config.GoEnv() == "dev" {
		ulid, err := algo.GenerateULID()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		fileName := fmt.Sprintf("%v_%v", ulid, file.Filename)
		dst := fmt.Sprintf("./static/img/profile/%v", fileName)

		// ローカルに画像を保存
		if err := c.SaveUploadedFile(file, dst); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		// 前の画像があれば削除
		if updateProfileImgResp.ProfileImgURL != "" {
			params := strings.Split(updateProfileImgResp.ProfileImgURL, "/")
			preFileName := params[len(params)-1]
			preFilePath := fmt.Sprintf("./static/img/profile/%v", preFileName)

			if err := os.Remove(preFilePath); err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		updateProfileImgResp.ProfileImgURL = fmt.Sprintf("%v/img/profile/%v", config.FilePath(), fileName)
	} else {
		// S3、またはCloud Storageでの保存処理

	}

	updateUserReq := &dto.UpdateUserRequest{
		Id:            id,
		ProfileImgURL: updateProfileImgResp.ProfileImgURL,
	}

	updateUserResp, err := uc.uu.Update(updateUserReq)
	if err != nil {
		if errors.Is(myerror.ErrDuplicatedKey, err) {
			c.JSON(http.StatusBadRequest, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, updateUserResp)
}
