package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dring1/jwt-oauth/middlewares"
	"github.com/dring1/jwt-oauth/models"
	"github.com/dring1/jwt-oauth/routes"
	"github.com/dring1/jwt-oauth/services"
	"github.com/justinas/alice"
)

func init() {
	services.Database()
	services.Database().HasTable(&models.User{})
}

func main() {
	router := routes.NewRouter()
	chain := alice.New(middlewares.LoggingHandler, middlewares.RecoverHandler).Then(router)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), chain)
	log.Fatal(err)
}
