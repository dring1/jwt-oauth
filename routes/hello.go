package routes

import (
	"github.com/dring1/jwt-oauth/controllers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func AddHelloRoutes(r *mux.Router) *mux.Router {
	r.Handle("/test/hello", alice.New().ThenFunc(controllers.HelloController)).Methods("Get")
	return r
}
