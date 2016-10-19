package routes

import (
	"net/http"

	"github.com/dring1/jwt-oauth/lib/contextkeys"
)

type ErrorRoute struct {
	Route
}

func (er *ErrorRoute) NewHandler() (*Route, error) {

	fn := func(w http.ResponseWriter, r *http.Request) {

		existingError := r.Context().Value(contextkeys.Error)
		if existingError == nil {
			existingError = "An error occurred"
		}

		w.WriteHeader(401)
		w.Write([]byte(existingError.(string)))
	}
	er.Handler = http.HandlerFunc(fn)
	return &er.Route, nil
}
