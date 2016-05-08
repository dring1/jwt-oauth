package routes

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r = AddHelloRoutes(r)
	r = LoginRoute(r)
	r = HomeRoute(r)
	return r
}
