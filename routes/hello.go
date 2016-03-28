package routes

import (
	"github.com/dring1/orm/controllers"
	"github.com/gorilla/mux"
)

func AddHelloRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/test/hello", controllers.HelloController).Methods("GET")
	return r
}
