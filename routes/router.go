package routes

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/dring1/jwt-oauth/middleware"
	"github.com/gorilla/mux"
)

// type Route interface {
// 	GenHttpHandlers() ([]*R, error)
// }
type Route struct {
	Path        string
	Methods     []string
	Middlewares []middleware.Middleware
	Handler     http.Handler
}

const (
	Get = "Get"
	Post = "POST"
)

// type R struct {
// 	Path    string
// 	Methods []string
// 	Handler http.Handler
// }

type RouteHandler interface {
	CompileRoute() (*Route, error)
}

type Router struct {
	*mux.Router
}

func New(gitHubClientID, gitHubClientSecret, redirectUrl string, services map[string]interface{}, middlewares map[string]middleware.Middleware) *Router {
	router := Router{mux.NewRouter()}
	routes := []RouteHandler{
		&GithubLoginRoute{Route: Route{Path: "/github/login", Methods: []string{Get}}, ClientID: gitHubClientID, ClientSecret: gitHubClientSecret},
		&GithubCallbackRoute{Route: Route{Path: "/github/callback", Methods: []string{Get}}, ClientID: gitHubClientID, ClientSecret: gitHubClientSecret},
		&UserProfileRoute{Route: Route{Path: "/profile", Methods: []string{Get}}},
		&HomeRoute{Route: Route{Path: "/", Methods: []string{Get}}, StaticFilePath: "./static"},
		&HelloRoute{Route: Route{Path: "/hello", Methods: []string{Get}, Middlewares: []middleware.Middleware{middlewares["VALIDATION"]}}},
		&TestRoute{Route: Route{Path: "/test", Methods: []string{Get}, Middlewares: []middleware.Middleware{middlewares["VALIDATION"]}}},
		&RefreshTokenRoute{Route: Route{Path: "/token/refresh", Methods: []string{Post}, Middlewares: []middleware.Middleware{middlewares["VALIDATION"]}}},
	}
	// Inject services into the routes | Tags or Reflection interface type impl
	// for _, r := range routes {
	// 	s := reflect.TypeOf(r).Elem()
	// 	for index := 0; index < s.NumField(); index++ {
	// 		field := s.Field(index)
	// 		if val, ok := field.Tag.Lookup("service"); ok && val != "" {
	// 			for _, ctrl := range controllers {
	// 				if reflect.TypeOf(ctrl).Elem().Name() == val {
	// 					val := reflect.ValueOf(r).Elem()
	// 					if val.Field(index).CanSet() {
	// 						val.Field(index).Set(reflect.ValueOf(ctrl))
	// 					}
	// 					break
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	routes = router.injectServices(routes, services)
	router.Register(routes)
	return &router
}

func (r *Router) Register(routes []RouteHandler) error {

	for _, route := range routes {
		h, err := route.CompileRoute()

		if err != nil {
			return err
		}
		handler := h.Handler
		if len(h.Middlewares) > 0 {
			handler = middleware.Handlers(h.Handler, h.Middlewares...)
		}
		r.Handle(h.Path, handler).Methods(h.Methods...)
		log.Printf("Registering %s with handlers for HTTP methods: %s", h.Path, strings.Join(h.Methods, ","))
	}
	return nil
}

func (r *Router) injectServices(routes []RouteHandler, services map[string]interface{}) []RouteHandler {
	for _, r := range routes {
		s := reflect.TypeOf(r).Elem()
		for index := 0; index < s.NumField(); index++ {
			field := s.Field(index)

			val, ok := field.Tag.Lookup("service")
			if !ok {
				continue
			}
			service, ok := services[val]
			log.Println("We are here", service, ok, val)
			if service, ok := services[val]; ok {
				val := reflect.ValueOf(r).Elem()
				if val.Field(index).CanSet() {
					val.Field(index).Set(reflect.ValueOf(service))
				}

			}
		}
	}
	return routes
}
