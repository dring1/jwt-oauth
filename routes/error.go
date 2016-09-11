package routes

import (
	"net/http"

	"github.com/dring1/jwt-oauth/lib/contextkeys"
)

type ErrorRoute struct {
	Route
}

func (er *ErrorRoute) NewHandler() (*R, error) {

	fn := func(w http.ResponseWriter, r *http.Request) {

		existingError := r.Context().Value(contextkeys.Error)
		if existingError == nil {
			existingError = "An error occurred"
		}

		w.WriteHeader(401)
		w.Write([]byte(existingError.(string)))
	}
	return &R{
		Path:    er.Path,
		Methods: er.Methods,
		Handler: http.HandlerFunc(fn),
	}, nil
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	existingError := r.Context().Value(contextkeys.Error)
	if existingError == nil {
		existingError = "An error occurred"
	}
	w.Write([]byte(existingError.(string)))
}
