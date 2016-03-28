package routes

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r = AddHelloRoutes(r)
	return r
}
