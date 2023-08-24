package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func userId(c *gin.Context) (string, error) {
	id, exists := c.Get("userId")
	if !exists {
		return "", errors.New("user id does not exist in context")
	}

	return id.(string), nil
}
