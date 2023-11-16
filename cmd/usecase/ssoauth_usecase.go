package usecase

import (
	"context"
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
	"english/config"
	"english/lib"
	"errors"
	"strings"
)

type SSOAuthUsecase interface {
	RedirectOAuthConsent(req *dto.RedirectOAuthConsentRequest) (*dto.RedirectOAuthConsentResponse, error)
	Callback(req *dto.CallbackRequest) (*dto.CallbackResponse, error)
}

type SSOAuthUsecaseImpl struct {
	ur   repository.UserRepository
	idPs map[string]*model.IdP
}

func NewSSOAuthUsecase(ur repository.UserRepository) SSOAuthUsecase {
	idPNames := []config.IdPName{config.GOOGLE, config.LINE}

	idPs := map[string]*model.IdP{}

	for _, name := range idPNames {
		idP := model.NewIdP(name)
		idPs[string(name)] = idP
	}

	return &SSOAuthUsecaseImpl{
		ur:   ur,
		idPs: idPs,
	}
}

func (lu *SSOAuthUsecaseImpl) RedirectOAuthConsent(req *dto.RedirectOAuthConsentRequest) (*dto.RedirectOAuthConsentResponse, error) {
	idPName := strings.ToUpper(req.IdPName)
	idP, ok := lu.idPs[idPName]
	if !ok {
		return nil, errors.New("invalid idP name")
	}

	state, err := lib.GenerateULID()
	if err != nil {
		return nil, err
	}

	redirectURL := idP.RedirectURL(state)

	resp := &dto.RedirectOAuthConsentResponse{
		RedirectURL: redirectURL,
		State:       state,
	}

	return resp, nil
}

func (lu *SSOAuthUsecaseImpl) Callback(req *dto.CallbackRequest) (*dto.CallbackResponse, error) {
	if req.QueryState != req.CookieState {
		return nil, errors.New("invalid state")
	}

	idpName := strings.ToUpper(req.IdpName)
	idP, ok := lu.idPs[idpName]
	if !ok {
		return nil, errors.New("invalid idP name")
	}

	idToken, err := idP.IdToken(context.Background(), req.Code)
	if err != nil {
		return nil, err
	}

	resp := &dto.CallbackResponse{}

	if err := idToken.Claims(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
