package routes

import (
	"net/http"
)

type ErrorRoute struct {
	R
}

func (er *ErrorRoute) NewHandler() (*R, error) {

	fn := func(w http.ResponseWriter, r *http.Request) {

		existingError := r.Context().Value(contextkeys.Error)
	}
	return &R{}, nil
}
