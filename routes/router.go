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
	Path        string
	Methods     []string
	Middlewares []http.Handler
}

type R struct {
	Path    string
	Methods []string
	Handler http.Handler
}

type RouteHandler interface {
	NewHandler() (*R, error)
}

type Router struct {
	*mux.Router
}

func New(gitHubClientID, gitHubClientSecret, redirectUrl string, controllers []controllers.Controller) *Router {
	router := Router{mux.NewRouter()}
	routes := []RouteHandler{
		&GithubLoginRoute{Route: Route{Path: "/github/login", Methods: []string{"GET"}}, ClientID: gitHubClientID, ClientSecret: gitHubClientSecret},
		&GithubCallbackRoute{Route: Route{Path: "/github/callback", Methods: []string{"GET"}}, ClientID: gitHubClientID, ClientSecret: gitHubClientSecret},
		&UserProfileRoute{Route: Route{Path: "/profile", Methods: []string{"GET"}}},
		&HomeRoute{Route: Route{Path: "/", Methods: []string{"GET"}}, StaticFilePath: "./static"},
		&HelloRoute{Route: Route{Path: "/hello", Methods: []string{"GET"}}},
		&ErrorRoute{Route: Route{Path: "/error", Methods: []string{"GET"}}},
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
		h, err := route.NewHandler()
		if err != nil {
			return err
		}
		r.Handle(h.Path, h.Handler).Methods(h.Methods...)
		log.Printf("Registering %s with handlers for HTTP methods: %s", h.Path, strings.Join(h.Methods, ","))
	}
	return nil
}
