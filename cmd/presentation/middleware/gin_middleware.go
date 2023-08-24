package middleware

import (
	"english/config"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type GinMiddleware struct {
}

func NewGinMiddleware() *GinMiddleware {
	return &GinMiddleware{}
}

func (gm *GinMiddleware) JWTMiddleware(c *gin.Context) {
	path := c.Request.URL.Path
	if path == "/" || path == "/login" || path == "/signup" || strings.Contains(path, "/auth") {
		c.Next()
		return
	}

	tokenStr, err := c.Cookie("token")
	if err != nil {
		log.Printf("Error: %v\n", err)
		c.JSON(http.StatusForbidden, "missing jwt token")
		c.Abort()
		return
	}

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method '%v'", token.Header["alg"])
		}

		return []byte(config.Secret()), nil
	}

	token, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		log.Printf("Error: %v\n", err)
		c.JSON(http.StatusForbidden, err.Error())
		c.Abort()
		return
	}

	if !token.Valid {
		message := "invalid token"
		log.Printf("Error: %v\n", message)
		c.JSON(http.StatusForbidden, message)
		c.Abort()
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	c.Set("userId", userId)

	c.Next()
}
