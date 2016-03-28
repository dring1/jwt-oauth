package main

import (
	"net/http"

	"github.com/dring1/orm/middlewares"
	"github.com/dring1/orm/routes"
	"github.com/justinas/alice"
)

func main() {
	router := routes.NewRouter()
	chain := alice.New(middlewares.LoggingHandler, middlewares.RecoverHandler).Then(router)
	http.ListenAndServe(":5000", chain)
}
