package main

import (
	"english/config"
	"english/src/infrastructure/dbconn"
	"english/src/infrastructure/repository"
	"english/src/presentation"
	"english/src/presentation/controller"
	"english/src/usecase/userusecase"
	"fmt"
	"io"
	"net/http"
)

func Top(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
}

func main() {
	db := dbconn.NewDB()
	ur := repository.NewUserMySQLRepository(db)
	ulu := userusecase.NewUserJWTLoginUsecase(ur)
	uc := controller.NewUserGinController(ulu)
	router := presentation.NewGinRouter(uc)

	router.Run(fmt.Sprintf(":%v", config.Port()))
}
