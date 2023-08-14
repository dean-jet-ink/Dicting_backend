package main

import (
	"io"
	"net/http"
)

func Top(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!!")
}

func main() {
	http.HandleFunc("/", Top)
	http.ListenAndServe(":8080", nil)
}
