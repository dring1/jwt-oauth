package routes

import (
	"log"
	"net/http"
	"reflect"
	"strings"

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

func New(gitHubClientID, gitHubClientSecret, redirectUrl string, services map[string]interface{}) *Router {
	router := Router{mux.NewRouter()}
	routes := []RouteHandler{
		&GithubLoginRoute{Route: Route{Path: "/github/login", Methods: []string{"GET"}}, ClientID: gitHubClientID, ClientSecret: gitHubClientSecret},
		&GithubCallbackRoute{Route: Route{Path: "/github/callback", Methods: []string{"GET"}}, ClientID: gitHubClientID, ClientSecret: gitHubClientSecret},
		&UserProfileRoute{Route: Route{Path: "/profile", Methods: []string{"GET"}}},
		&HomeRoute{Route: Route{Path: "/", Methods: []string{"GET"}}, StaticFilePath: "./static"},
		&HelloRoute{Route: Route{Path: "/hello", Methods: []string{"GET"}}},
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
		h, err := route.NewHandler()
		if err != nil {
			return err
		}
		r.Handle(h.Path, h.Handler).Methods(h.Methods...)
		log.Printf("Registering %s with handlers for HTTP methods: %s", h.Path, strings.Join(h.Methods, ","))
	}
	return nil
}

func (r *Router) injectServices(routes []RouteHandler, services map[string]interface{}) []RouteHandler {

	// for _, route := range routes {
	// 	val := reflect.ValueOf(route).Elem()
	// 	for index := 0; index < val.NumField(); index++ {
	// 		for _, s := range services {
	// 			field := val.Type().Field(index).Type

	// 			// f := field.Elem()

	// 			// n := reflect.New(field).Interface()

	// 			x = reflect.TypeOf(s)

	// 			log.Println(field, x)

	// 			// fieldInterface := reflect.TypeOf(f)
	// 			// serviceValue := reflect.ValueOf(s).Elem()

	// 			// serviceType := serviceValue.Type()
	// 			// serviceType.Typ

	// 			// fieldInterfaceType := reflect.TypeOf(val.Type().Field(index))
	// 			// x := reflect.TypeOf(fieldType).Elem()
	// 			// ms := reflect.New(x).Elem().Interface()
	// 			// log.Println(x, ms)
	// 			// t := serviceType.Convert(fieldType)
	// 			// fmt.Println(t)

	// 			// serviceType.(fieldType)
	// 			// (reflect.ValueOf(fieldType).Type())(serviceType)
	// 			// log.Println(f, serviceType, fieldInterface)
	// 			// ok := serviceType.Implements(f)
	// 			// if f == reflect.TypeOf(serviceType) {
	// 			// 	if x := val.Field(index); x.CanSet() {
	// 			// 		log.Println("Hi")
	// 			// 		x.Set(reflect.ValueOf(s))
	// 			// 	}
	// 			// }
	// 		}
	// 	}
	// }
	// Inject services into the routes | Tags or Reflection interface type impl
	log.Println(services)

	for _, r := range routes {
		s := reflect.TypeOf(r).Elem()
		for index := 0; index < s.NumField(); index++ {
			field := s.Field(index)

			val, ok := field.Tag.Lookup("service")
			log.Printf("We got our tags %v %s", val, ok)
			if !ok {
				continue
			}
			service, ok := services[val]
			log.Println("We are here", service, ok, val)
			if service, ok := services[val]; ok {
				val := reflect.ValueOf(r).Elem()
				if val.Field(index).CanSet() {
					log.Println("Setting")
					val.Field(index).Set(reflect.ValueOf(service))
				}

			}

			// for key, service := range services {
			// 	if val, ok := field.Tag.Lookup("service"); ok && val != "" {
			// 		// if reflect.TypeOf(ctrl).Elem().Name() == val {
			// 		if key != val {
			// 			continue
			// 		}
			// 		val := reflect.ValueOf(r).Elem()
			// 		if val.Field(index).CanSet() {
			// 			log.Println("Setting")
			// 			val.Field(index).Set(reflect.ValueOf(service))
			// 		}
			// 		break

			// 	}
			// }
		}
	}
	return routes
}
