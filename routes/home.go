package routes

import (
	"net/http"
	"path/filepath"

	"github.com/dring1/jwt-oauth/controllers"
)

type HomeRoute struct {
	Route
	StaticFilePath string
	Controller     controllers.Controller `controller:"HelloController"`
}

func (r *HomeRoute) ServeHTTP(w http.ResponseWriter, res *http.Request) {
	_, err := filepath.Abs(r.StaticFilePath)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error"))
	}
	http.StripPrefix("/", http.FileServer(http.Dir(r.StaticFilePath)))
}

func (r *HomeRoute) GetPath() string {
	return r.Route.Path
}

func (r *HomeRoute) GetMethods() []string {
	return r.Route.Methods
}
