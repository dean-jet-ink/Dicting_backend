package usecase

import (
	"context"
	"english/algo"
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/usecase/dto"
	"errors"
)

type SSOLoginUsecase interface {
	RedirectOAuthConsent(req *dto.RedirectOAuthConsentRequest) (*dto.RedirectOAuthConsentResponse, error)
	Callback(req *dto.CallbackRequest) (*dto.CallbackResponse, error)
}

type OIDCLoginUsecase struct {
	ur   repository.UserRepository
	idPs map[string]*model.IdP
}

func NewOIDCLoginUsecase(ur repository.UserRepository) SSOLoginUsecase {
	idPNames := []string{"google", "line"}

	idPs := map[string]*model.IdP{}

	for _, name := range idPNames {
		idP := model.NewIdP(name)
		idPs[name] = idP
	}

	return &OIDCLoginUsecase{
		ur:   ur,
		idPs: idPs,
	}
}

func (lu *OIDCLoginUsecase) RedirectOAuthConsent(req *dto.RedirectOAuthConsentRequest) (*dto.RedirectOAuthConsentResponse, error) {
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

	output := &dto.RedirectOAuthConsentResponse{
		RedirectURL: redirectURL,
		State:       state,
	}

	return output, nil
}

func (lu *OIDCLoginUsecase) Callback(req *dto.CallbackRequest) (*dto.CallbackResponse, error) {
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

	output := &dto.CallbackResponse{}

	if err := idToken.Claims(output); err != nil {
		return nil, err
	}

	return output, nil
}
