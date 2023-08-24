package router

import (
	"english/cmd/presentation/controller"
	"english/cmd/presentation/middleware"
	"text/template"

	"github.com/gin-gonic/gin"
)

func NewGinRouter(uc controller.UserController) *gin.Engine {
	router := gin.Default()
	mid := middleware.NewGinMiddleware()

	router.Use(mid.JWTMiddleware)

	router.GET("/", func(c *gin.Context) {
		templ, _ := template.ParseFiles("mock/sso_test.html")
		templ.Execute(c.Writer, nil)
	})

	// 認証関係
	router.POST("/signup", uc.Signup)
	router.POST("/login", uc.Login)
	router.POST("/logout", uc.Logout)
	router.GET("/auth", uc.RedirectOAuthConsent)
	router.GET("/auth/callback", uc.OAuthCallback)

	router.POST("/user/update", uc.UpdateProfile)

	return router
}
