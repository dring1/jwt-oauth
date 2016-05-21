package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dring1/orm/middlewares"
	"github.com/dring1/orm/models"
	"github.com/dring1/orm/routes"
	"github.com/dring1/orm/services"
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
