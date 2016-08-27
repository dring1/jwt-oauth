package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// type Route interface {
// 	GenHttpHandlers() ([]*R, error)
// }
type R struct {
	Route string
	http.Handler
	Methods []string
}

type Route interface {
	http.Handler
	GetRoute() string
	GetMethods() []string
}

func New(gitHubClientID, gitHubClientSecret string) *mux.Router {
	router := mux.NewRouter()
	routes := []Route{
		// &LoginRoute{
		// 	GitHubClientID:     gitHubClientID,
		// 	GitHubClientSecret: gitHubClientSecret,
		// },
		// &HelloRoute{},
		&HomeRoute{StaticFilePath: "static"},
	}
	register(router, routes)
	return router
}

func register(router *mux.Router, routes []Route) error {

	for _, route := range routes {
		router.Handle(route.GetRoute(), route)
		// rs, err := route.GenHttpHandlers()
		// if err != nil {
		// 	return err
		// }
		// for _, r := range rs {
		// 	log.Printf("Registering %s with handlers for HTTP methods: %s", r.Path, strings.Join(r.Methods, ","))
		// 	router.Handle(r.Path, r.Handler).Methods(r.Methods...)
		// }
	}
	return nil
}
