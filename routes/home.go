package routes

import (
	"net/http"
	"path/filepath"
)

type HomeRoute struct {
	Route
	StaticFilePath string
}

func (r *HomeRoute) CompileRoute() (*Route, error) {
	_, err := filepath.Abs(r.StaticFilePath)
	if err != nil {
		return nil, err
	}
	r.Handler = http.StripPrefix("/", http.FileServer(http.Dir(r.StaticFilePath)))
	return &r.Route, nil
}
