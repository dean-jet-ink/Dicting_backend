package router

import (
	"english/cmd/presentation/controller"
	"english/cmd/presentation/middleware"
	"english/config"
	"text/template"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewGinRouter(uc controller.UserController, ec controller.EnglishItemController, oc controller.OutputController) *gin.Engine {
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
	router.PUT("/user", uc.UpdateProfile)
	router.PUT("/user/profile-img", uc.UpdateProfileImg)

	router.GET("/english", ec.GetByUserId)
	router.GET("/english/:id", ec.GetById)
	router.GET("/english/proposal", ec.Proposal)
	router.GET("/english/proposal/translation", ec.ProposalTranslation)
	router.GET("/english/proposal/explanation", ec.ProposalExplanation)
	router.GET("/english/proposal/example", ec.ProposalExample)
	router.GET("/english/required-exp", ec.GetRequiredExp)
	router.POST("/english", ec.Create)
	router.PUT("/english", ec.Update)
	router.DELETE("/english/:id", ec.Delete)

	router.GET("/english/output", oc.GetQuestion)
	router.POST("/english/output", oc.AnswerQuestions)
	router.PUT("/english/output", oc.CreateOutput)
	router.GET("/english/output/times/:englishItemId", oc.GetOutputTimes)
	router.GET("/english/outputs", oc.GetOutputs)
	router.DELETE("/english/output", oc.DeleteOutput)

	return router
}
