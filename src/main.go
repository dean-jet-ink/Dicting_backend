package main

import (
	"english/config"
	"io"
	"log"
	"net/http"
)

func Top(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
}

func main() {
	config.SetLogger(config.LogFileName())

	log.Println(config.GoEnv())
	log.Println(config.Port())
	log.Println(config.APIDomain())
	log.Println(config.FrontEndURL())
	log.Println(config.Secret())
	log.Println(config.LogFileName())
	log.Println(config.MySQLDBName())
	log.Println(config.MySQLUser())
	log.Println(config.MySQLPass())
	log.Println(config.MySQLHost())

	http.HandleFunc("/", Top)
	http.ListenAndServe(":8080", nil)
}
