package routes

import (
	"context"
	"net/http"

	"github.com/dring1/jwt-oauth/lib/contextkeys"
	"github.com/dring1/jwt-oauth/lib/errors"
)

type HelloRoute struct {
	Route
	StaticFilePath string
	// Controller     controllers.Controller `controller:"HelloController"`
}

func (r *HelloRoute) NewHandler() (*R, error) {

	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextkeys.Error, errors.UnauthorizedUser)
		r = r.WithContext(ctx)
		w.WriteHeader(401)
		ErrorHandler(w, r)
	}
	return &R{
		Path:    r.Path,
		Methods: r.Methods,
		Handler: http.HandlerFunc(fn),
	}, nil
}
