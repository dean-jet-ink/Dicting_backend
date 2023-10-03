package router

import (
	"english/cmd/presentation/controller"
	"english/cmd/presentation/middleware"
	"english/config"
	"text/template"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewGinRouter(uc controller.UserController, ec controller.EnglishItemController) *gin.Engine {
	router := gin.New()
	mid := middleware.NewGinMiddleware()

	// cors制約の設定
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			config.FrontEndURL(),
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 静的ファイルの設定
	router.Static("/static", "./static")

	// パニック時のミドルウェア
	router.Use(mid.RecoverPanic)

	// jwtの有効性の確認、及びjwt内のuser idをcontextに格納
	router.Use(mid.JWTMiddleware)

	router.GET("/", func(c *gin.Context) {
		templ, _ := template.ParseFiles("mock/sso_test.html")
		templ.Execute(c.Writer, nil)
	})

	router.POST("/signup", uc.Signup)
	router.POST("/login", uc.Login)
	router.POST("/logout", uc.Logout)
	router.GET("/auth", uc.RedirectOAuthConsent)
	router.GET("/auth/callback", uc.OAuthCallback)

	router.GET("/user", uc.GetUser)
	router.POST("/user/update", uc.UpdateProfile)
	router.POST("/user/update/profile-img", uc.UpdateProfileImg)

	router.GET("/english", ec.GetByUserIdAndContent)
	router.GET("/english/proposal", ec.Proposal)
	router.POST("/english", ec.Create)

	return router
}
