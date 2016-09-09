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

func (r *HomeRoute) NewHandler() (*R, error) {
	_, err := filepath.Abs(r.StaticFilePath)
	if err != nil {
		return nil, err
	}

	return &R{
		Path:    r.Path,
		Methods: r.Methods,
		Handler: http.StripPrefix("/", http.FileServer(http.Dir(r.StaticFilePath))),
	}, nil
}
