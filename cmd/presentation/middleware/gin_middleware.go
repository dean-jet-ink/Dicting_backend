package middleware

import (
	"english/cmd/presentation/errhandle"
	"english/config"
	"english/myerror"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type GinMiddleware struct {
}

func NewGinMiddleware() *GinMiddleware {
	return &GinMiddleware{}
}

func (gm *GinMiddleware) RecoverPanic(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			errhandle.HandleErrorJSON(err, c)
		}
	}()
	c.Next()
}

func (gm *GinMiddleware) JWTMiddleware(c *gin.Context) {
	path := c.Request.URL.Path
	if path == "/" || path == "/login" || path == "/signup" || strings.Contains(path, "/auth") {
		c.Next()
		return
	}

	tokenStr, err := c.Cookie("token")
	if err != nil {
		errhandle.HandleErrorJSON(myerror.ErrMissingJWT, c)
		c.Abort()
		return
	}

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%v '%w'", myerror.ErrUnexpectedSigningMethod, errors.New(token.Header["alg"].(string)))
		}

		return []byte(config.JWTSecret()), nil
	}

	token, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		errhandle.HandleErrorJSON(err, c)
		c.Abort()
		return
	}

	if !token.Valid {
		errhandle.HandleErrorJSON(myerror.ErrInvalidToken, c)
		c.Abort()
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	c.Set("userId", userId)

	c.Next()
}
