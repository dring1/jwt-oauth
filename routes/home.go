package routes

import (
	"net/http"
	"path/filepath"
)

type HomeRoute struct {
	R
	StaticFilePath string
}

func (r *HomeRoute) ServeHTTP(w http.ResponseWriter, res *http.Request) {
	_, err := filepath.Abs(r.StaticFilePath)
	http.StripPrefix("/", http.FileServer(http.Dir(r.StaticFilePath)))
}

func (r *HomeRoute) GetRoute() string {
	return r.Route
}

func (r *HomeRoute) GetMethods() []string {
	return r.Methods
}
