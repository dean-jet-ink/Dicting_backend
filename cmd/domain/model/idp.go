package model

import (
	"context"
	"english/config"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2"
)

var scopes = []string{oidc.ScopeOpenID, "email", "profile"}

type IdP struct {
	idPName      string
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func NewIdP(idPName string) *IdP {
	provider, _ := oidc.NewProvider(context.Background(), config.ISSURL(idPName))
	clientID := config.ClientId(idPName)

	oauth2config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: config.ClientSecret(idPName),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  config.OAuthRedirectURL(),
		Scopes:       scopes,
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	verifier := provider.Verifier(oidcConfig)

	return &IdP{
		oauth2Config: oauth2config,
		verifier:     verifier,
	}
}

func (idP *IdP) RedirectURL(state string) string {
	return idP.oauth2Config.AuthCodeURL(state)
}

func (idP *IdP) IdToken(ctx context.Context, code string) (*oidc.IDToken, error) {
	oauthToken, err := idP.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	rawIDToken := oauthToken.Extra("id_token").(string)

	idToken, err := idP.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		if matched := strings.Contains(err.Error(), "HS256"); matched {
			verifier := idP.createHS256Verifier()

			idToken, err = verifier.Verify(ctx, rawIDToken)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return idToken, nil
}

// go-oidcがHS256に非対応のため、中でjoseがHS256に対応しているので
// それをここでgo-oidcの形で再定義し、go-oidcに渡す
type HS256KeySet struct {
	CommonKey string
}

func (h *HS256KeySet) verify(_ context.Context, jws *jose.JSONWebSignature) ([]byte, error) {
	return jws.Verify([]byte(h.CommonKey))
}

func (h *HS256KeySet) VerifySignature(ctx context.Context, jwt string) ([]byte, error) {
	jws, err := jose.ParseSigned(jwt)
	if err != nil {
		return nil, fmt.Errorf("oidc: malformed jwt: %v", err)
	}

	return h.verify(ctx, jws)
}

func (idP *IdP) createHS256Verifier() *oidc.IDTokenVerifier {
	keySet := &HS256KeySet{
		CommonKey: config.ClientSecret(idP.idPName),
	}

	oicdConfig := &oidc.Config{
		ClientID:             config.ClientId(idP.idPName),
		SupportedSigningAlgs: []string{"HS256"},
	}

	verifier := oidc.NewVerifier(config.ISSURL(idP.idPName), keySet, oicdConfig)

	return verifier
}
