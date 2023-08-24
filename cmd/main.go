package main

import (
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
	su := usecase.NewStandardSignupUsecase(ur)
	lu := usecase.NewStandardLoginUsecase(ur)
	ssu := usecase.NewOIDCAuthUsecase(ur)
	uu := usecase.NewUpdateUserProfileUsecase(ur)
	uc := controller.NewUserGinController(su, lu, ssu, uu)
	router := router.NewGinRouter(uc)

	router.Run(fmt.Sprintf(":%v", config.Port()))
}
