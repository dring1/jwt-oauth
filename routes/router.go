package routes

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/dring1/jwt-oauth/middleware"
	"github.com/dring1/jwt-oauth/services"
	"github.com/gorilla/mux"
)

type Responder http.Handler

type Route struct {
	Path        string
	Methods     []string
	Middlewares []middleware.Middleware
	Handler     http.Handler
	Respond     Responder
}

type RouteRaw interface {
	CompileRoute() (*Route, error)
}

type Router struct {
	*mux.Router
}

type Config struct {
	Services     *services.Services
	Middlewares  middleware.MiddlewareMap
	ClientID     string
	ClientSecret string
}

const (
	Get  = "GET"
	Post = "POST"
)

func NewRouter(routes []*Route) (*Router, error) {
	router := Router{mux.NewRouter()}
	err := router.Register(routes)
	if err != nil {
		return nil, err
	}
	return &router, nil
}

func NewRoutes(config *Config) ([]*Route, error) {
	routes := []RouteRaw{
		&GithubLoginRoute{
			Route: Route{
				Path:    "/github/login",
				Methods: []string{Get},
			},
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
		},
		&GithubCallbackRoute{
			Route: Route{
				Path:    "/github/callback",
				Methods: []string{Get},
			},
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
		},
		&HomeRoute{Route: Route{Path: "/", Methods: []string{Get}}, StaticFilePath: "./static"},
		&TestRoute{Route: Route{Path: "/test", Methods: []string{Get}, Middlewares: []middleware.Middleware{config.Middlewares[middleware.ValidateMiddleware]}}},
		&RefreshTokenRoute{Route: Route{Path: "/token/refresh", Methods: []string{Get}, Middlewares: []middleware.Middleware{config.Middlewares[middleware.ValidateMiddleware]}}},
	}
	hydratedRoutes := InjectServices(routes, config.Services)

	r, err := TransformRoutes(hydratedRoutes)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func TransformRoutes(routesRaw []RouteRaw) ([]*Route, error) {
	routes := []*Route{}
	for _, route := range routesRaw {
		r, err := route.CompileRoute()
		if err != nil {
			return nil, err
		}
		r.Handler = middleware.Handlers(r.Handler, r.Middlewares...)
		routes = append(routes, r)
	}
	return routes, nil
}

func (r *Router) Register(routes []*Route) error {
	for _, route := range routes {
		r.Handle(route.Path, route.Handler).Methods(route.Methods...)
		log.Printf("Registering %s with handlers for HTTP methods: %s", route.Path, strings.Join(route.Methods, ","))
	}
	return nil
}

func InjectServices(routes []RouteRaw, svcs *services.Services) []RouteRaw {
	sv := reflect.ValueOf(svcs).Elem()
	for _, r := range routes {
		s := reflect.TypeOf(r).Elem()
		for index := 0; index < s.NumField(); index++ {
			field := s.Field(index)
			tagValue, ok := field.Tag.Lookup("service")
			if !ok {
				continue
			}
			serviceValue := sv.FieldByName(tagValue).Elem()
			val := reflect.ValueOf(r).Elem()
			if val.Field(index).CanSet() {
				val.Field(index).Set(serviceValue)
			}
		}
	}
	return routes
}
