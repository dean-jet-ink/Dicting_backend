package main

import (
	"english/cmd/infrastructure/client"
	"english/cmd/infrastructure/dbconn"
	"english/cmd/infrastructure/gateway"
	"english/cmd/presentation/controller"
	"english/cmd/presentation/router"
	"english/cmd/usecase"
	"english/config"
	"fmt"
)

func init() {
	config.SetLogger(config.LogFileName())
}

func main() {
	db := dbconn.NewDB()

	userRepo := gateway.NewUserMySQLRepository(db)
	signupUse := usecase.NewSignupUsecase(userRepo)
	loginUse := usecase.NewLoginUsecase(userRepo)
	ssoAuthUse := usecase.NewSSOAuthUsecase(userRepo)
	updateUserUse := usecase.NewUpdateUserProfileUsecase(userRepo)
	updateProfileUse := usecase.NewUpdateProfileImgUsecase(userRepo)
	userGinCon := controller.NewUserGinController(signupUse, loginUse, ssoAuthUse, updateUserUse, updateProfileUse)

	chatGPTAPI := client.NewOpenAIAPI()
	proposalUse := usecase.NewProposalEnglishItemUsecase(chatGPTAPI)
	englishItemRepo := gateway.NewEnglishItemMySQLReporitory(db)
	fileStorageRepo := gateway.NewFileStorageGCSRepository()
	createEnglishItemUse := usecase.NewCreateEnglishItemUsecase(englishItemRepo, fileStorageRepo)
	ec := controller.NewEnglishItemController(proposalUse, createEnglishItemUse)

	router := router.NewGinRouter(userGinCon, ec)

	router.Run(fmt.Sprintf(":%v", config.Port()))
}
