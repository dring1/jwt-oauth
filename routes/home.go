package routes

import (
	"net/http"
	"path/filepath"
)

type HomeRoute struct {
	Route
	StaticFilePath string
	// Controller     controllers.Controller `controller:"HelloController"`
}

func (r *HomeRoute) ServeHTTP(w http.ResponseWriter, res *http.Request) {
	_, err := filepath.Abs(r.StaticFilePath)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
	}
	http.StripPrefix("/", http.FileServer(http.Dir(r.StaticFilePath))).ServeHTTP(w, res)
}

func (r *HomeRoute) GetPath() string {
	return r.Route.Path
}

func (r *HomeRoute) GetMethods() []string {
	return r.Route.Methods
}
