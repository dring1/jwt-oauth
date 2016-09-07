package routes

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/dring1/jwt-oauth/controllers"
	"github.com/gorilla/mux"
)

// type Route interface {
// 	GenHttpHandlers() ([]*R, error)
// }
type Route struct {
	http.Handler
	Path    string
	Methods []string
}

type RouteHandler interface {
	http.Handler
	GetPath() string
	GetMethods() []string
}

type Router struct {
	*mux.Router
}

func New(gitHubClientID, gitHubClientSecret string, controllers []controllers.Controller) *Router {
	router := Router{mux.NewRouter()}
	routes := []RouteHandler{
		// &LoginRoute{
		// 	GitHubClientID:     gitHubClientID,
		// 	GitHubClientSecret: gitHubClientSecret,
		// },
		// &HelloRoute{},
		&HomeRoute{Route: Route{Path: "/home", Methods: []string{"GET"}}, StaticFilePath: "static"},
	}
	for _, r := range routes {
		s := reflect.TypeOf(r).Elem()
		for index := 0; index < s.NumField(); index++ {
			field := s.Field(index)
			if val, ok := field.Tag.Lookup("controller"); ok && val != "" {
				for _, ctrl := range controllers {
					if reflect.TypeOf(ctrl).Elem().Name() == val {
						val := reflect.ValueOf(r).Elem()
						if val.Field(index).CanSet() {
							val.Field(index).Set(reflect.ValueOf(ctrl))
						}
						break
					}

				}
			}

		}

	}
	router.Register(routes)
	return &router
}

func (r *Router) Register(routes []RouteHandler) error {

	for _, route := range routes {
		r.Handle(route.GetPath(), route).Methods(route.GetMethods()...)
		log.Printf("Registering %s with handlers for HTTP methods: %s", route.GetPath(), strings.Join(route.GetMethods(), ","))
	}
	return nil
}
