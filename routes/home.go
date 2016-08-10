package routes

import (
	"net/http"
	"path/filepath"
)

type HomeRoute struct {
	StaticFilePath string
}

func (r *HomeRoute) GenHttpHandlers() ([]*R, error) {
	_, err := filepath.Abs(r.StaticFilePath)
	if err != nil {
		return nil, err
	}
	return []*R{
		&R{
			Path:    "/",
			Methods: []string{"GET"},
			Handler: http.StripPrefix("/", http.FileServer(http.Dir(r.StaticFilePath))),
		},
	}, nil
}
