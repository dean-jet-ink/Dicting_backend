package dto

type RedirectOAuthConsentRequest struct {
	IdPName string `form:"idp_name" validate:"required"`
}

type RedirectOAuthConsentResponse struct {
	RedirectURL string
	State       string
}

type CallbackRequest struct {
	IdpName     string
	QueryState  string `form:"state" validate:"required"`
	CookieState string
	Code        string `form:"code" validate:"required"`
}

type CallbackResponse struct {
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
}