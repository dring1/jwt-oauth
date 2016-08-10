package routes

import (
	"net/http"

	"github.com/dring1/jwt-oauth/controllers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func AddHelloRoutes(r *mux.Router) *mux.Router {
	r.Handle("/test/hello", alice.New().ThenFunc(controllers.HelloController)).Methods("Get")
	return r
}

type HelloRoute struct{}

func (h *HelloRoute) GenHttpHandlers() ([]*R, error) {
	x := alice.New().ThenFunc(controllers.HelloController).(http.HandlerFunc)
	return []*R{&R{Path: "/test/hello", Handler: x, Methods: []string{"GET"}}}, nil
}
