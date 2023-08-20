package presentation

import (
	"english/src/presentation/controller"

	"github.com/gin-gonic/gin"
)

func NewGinRouter(uc controller.UserController) *gin.Engine {
	router := gin.Default()

	router.POST("/login", uc.Login)
	router.POST("/logout", uc.Logout)
	router.GET("/auth", uc.RedirectOAuthConsent)
	router.GET("/auth/callback", uc.OAuthCallback)

	return router
}
