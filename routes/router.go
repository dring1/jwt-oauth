package routes

import "github.com/gorilla/mux"

type RouterConfig struct {
	GitHubID string
}

func NewRouter(c *RouterConfig) *mux.Router {
	r := mux.NewRouter()
	r = AddHelloRoutes(r)
	r = LoginRoute(r, nil, nil, nil)
	r = HomeRoute(r)
	return r
}
