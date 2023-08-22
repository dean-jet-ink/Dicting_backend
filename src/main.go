package main

import (
	"english/config"
	"english/src/infrastructure/dbconn"
	"english/src/infrastructure/repository"
	"english/src/presentation"
	"english/src/presentation/controller"
	"english/src/usecase"
	"fmt"
	"io"
	"net/http"
)

func Top(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
}

func main() {
	fmt.Println("hello")
	db := dbconn.NewDB()
	ur := repository.NewUserMySQLRepository(db)
	lu := usecase.NewJWTLoginUsecase(ur)
	su := usecase.NewOIDCLoginUsecase(ur)
	uc := controller.NewUserGinController(lu, su)
	router := presentation.NewGinRouter(uc)

	router.Run(fmt.Sprintf(":%v", config.Port()))
}
