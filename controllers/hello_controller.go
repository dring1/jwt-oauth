package controllers

import (
	"log"
	"net/http"
)

func HelloController(w http.ResponseWriter, r *http.Request) {
	log.Println("Hi!")
	w.Write([]byte("<h1>Hello, world!</h1>"))
}
