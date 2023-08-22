package presentation

import (
	"english/cmd/presentation/controller"

	"github.com/gin-gonic/gin"
)

func NewGinRouter(uc controller.UserController) *gin.Engine {
	router := gin.Default()

	router.POST("/signup", uc.Signup)
	router.POST("/login", uc.Login)
	router.POST("/logout", uc.Logout)

	return router
}
