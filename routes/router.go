package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type RouterConfig struct {
	GitHubID string
}

func NewRouter(gitHubClientID, githubClientSecret string) *mux.Router {
	r := mux.NewRouter()
	r = AddHelloRoutes(r)
	r = LoginRoute(r, nil, nil, nil)
	r = HomeRoute(r)
	return r
}

func New() *mux.Router {
	routes := []Route{
		LoginRoute{
			GitHubClientID:     gitHubClientID,
			GitHubClientSecret: gitHubClientSecret,
		},
	}
	return nil
}

type Route interface {
	HttpHandler() map[string]http.HandlerFunc
}
