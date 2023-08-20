package usecase

import (
	"context"
	"english/algo"
	"english/src/domain/usermodel"
	"errors"
)

type SSOLoginUsecase interface {
	RedirectOAuthConsent(req *RedirectOAuthConsentRequest) (*RedirectOAuthConsentOutput, error)
	Callback(req *CallbackRequest) (*CallbackOutput, error)
}

type OIDCLoginUsecase struct {
	ur   usermodel.UserRepository
	idPs map[string]*usermodel.IdP
}

func NewOIDCLoginUsecase(ur usermodel.UserRepository) SSOLoginUsecase {
	idPNames := []string{"google", "line"}

	idPs := map[string]*usermodel.IdP{}

	for _, name := range idPNames {
		idP := usermodel.NewIdP(name)
		idPs[name] = idP
	}

	return &OIDCLoginUsecase{
		ur:   ur,
		idPs: idPs,
	}
}

func (lu *OIDCLoginUsecase) RedirectOAuthConsent(req *RedirectOAuthConsentRequest) (*RedirectOAuthConsentOutput, error) {
	idPName := req.IdPName
	idP, ok := lu.idPs[idPName]
	if !ok {
		return nil, errors.New("invalid idP name")
	}

	state, err := algo.GenerateULID()
	if err != nil {
		return nil, err
	}

	redirectURL := idP.RedirectURL(state)

	output := &RedirectOAuthConsentOutput{
		redirectURL: redirectURL,
		state:       state,
	}

	return output, nil
}

func (lu *OIDCLoginUsecase) Callback(req *CallbackRequest) (*CallbackOutput, error) {
	if req.QueryState != req.CookieState {
		return nil, errors.New("invalid state")
	}

	idP, ok := lu.idPs[req.IdpName]
	if !ok {
		return nil, errors.New("invalid idP name")
	}

	idToken, err := idP.IdToken(context.Background(), req.Code)
	if err != nil {
		return nil, err
	}

	output := &CallbackOutput{}

	if err := idToken.Claims(output); err != nil {
		return nil, err
	}

	return output, nil
}

type RedirectOAuthConsentRequest struct {
	IdPName string `json:"idp_name" validate:"required"`
}

type RedirectOAuthConsentOutput struct {
	redirectURL string
	state       string
}

func (ro *RedirectOAuthConsentOutput) RedirectURL() string {
	return ro.redirectURL
}

func (ro *RedirectOAuthConsentOutput) State() string {
	return ro.state
}

type CallbackRequest struct {
	IdpName     string
	QueryState  string `json:"state" validate:"required"`
	CookieState string
	Code        string `json:"code" validate:"required"`
}

type CallbackOutput struct {
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
}
