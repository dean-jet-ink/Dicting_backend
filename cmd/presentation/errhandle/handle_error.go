package errhandle

import (
	"english/myerror"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func handleError(err error, c *gin.Context, fn func(code int, resp *ErrorResponse, c *gin.Context)) {
	log.Println(err)

	resp := &ErrorResponse{Message: err.Error()}

	if errors.Is(myerror.ErrDuplicatedKey, err) ||
		errors.Is(myerror.ErrRecordNotFound, err) ||
		errors.Is(myerror.ErrValidation, err) ||
		errors.Is(myerror.ErrBindingFailure, err) ||
		errors.Is(myerror.ErrValidation, err) ||
		errors.Is(myerror.ErrCookieExpired, err) ||
		errors.Is(myerror.ErrParsingFailure, err) ||
		errors.Is(myerror.ErrFormFileNotFound, err) {
		// 400
		fn(http.StatusBadRequest, resp, c)
	} else if errors.Is(myerror.ErrMismatchedPassword, err) {
		// 401
		fn(http.StatusUnauthorized, resp, c)
	} else if errors.Is(myerror.ErrUnverifiedEmail, err) ||
		errors.Is(myerror.ErrMissingJWT, err) ||
		errors.Is(myerror.ErrInvalidToken, err) ||
		errors.Is(myerror.ErrUnexpectedSigningMethod, err) {
		// 403
		fn(http.StatusForbidden, resp, c)
	} else {
		// 500
		fn(http.StatusInternalServerError, resp, c)
	}
}

func HandleErrorJSON(err error, c *gin.Context) {
	handleError(err, c, handleErrorJSON)
}

func handleErrorJSON(code int, resp *ErrorResponse, c *gin.Context) {
	c.JSON(code, resp)
}

func HandleValidationError(err error, c *gin.Context) {
	validationErr := fmt.Errorf("%v: %w", myerror.ErrValidation, err)
	handleError(validationErr, c, handleErrorJSON)
}

func HandleErrorRedirect(err error, c *gin.Context, location string) {
	fn := func(_ int, resp *ErrorResponse, c *gin.Context) {
		redirectURL := fmt.Sprintf("%v?error=%s", location, resp.Message)
		parsedURL, _ := url.Parse(redirectURL)

		c.Redirect(http.StatusFound, parsedURL.String())
	}

	handleError(err, c, fn)
}
