package main

import (
	"english/cmd/infrastructure/dbconn"
	"english/cmd/infrastructure/gateway"
	"english/cmd/presentation"
	"english/cmd/presentation/controller"
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
	uc := controller.NewUserGinController(su, lu, ssu)
	router := presentation.NewGinRouter(uc)

	router.Run(fmt.Sprintf(":%v", config.Port()))
}
