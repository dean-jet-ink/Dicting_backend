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

	ur := gateway.NewUserMySQLRepository(db)
	su := usecase.NewSignupUsecase(ur)
	lu := usecase.NewLoginUsecase(ur)
	ssu := usecase.NewSSOAuthUsecase(ur)
	uu := usecase.NewUpdateUserProfileUsecase(ur)
	upu := usecase.NewUpdateProfileImgUsecase(ur)
	uc := controller.NewUserGinController(su, lu, ssu, uu, upu)

	chatGPTAPI := client.NewOpenAIAPI()
	pu := usecase.NewProposalEnglishItemUsecase(chatGPTAPI)
	ec := controller.NewEnglishItemController(pu)

	router := router.NewGinRouter(uc, ec)

	router.Run(fmt.Sprintf(":%v", config.Port()))
}
