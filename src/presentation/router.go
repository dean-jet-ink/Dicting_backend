package presentation

import (
	"english/src/presentation/controller"

	"github.com/gin-gonic/gin"
)

func NewGinRouter(uc controller.UserController) *gin.Engine {
	router := gin.Default()

	router.POST("/login", uc.Login)
	router.POST("/logout", uc.Logout)

	return router
}
