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
}

func (r *HelloRoute) CompileRoute() (*Route, error) {

	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextkeys.Error, errors.UnauthorizedUser)
		r = r.WithContext(ctx)
		w.WriteHeader(401)
		errors.ErrorHandler(w, r)
	}
	r.Handler = http.HandlerFunc(fn)
	return &r.Route, nil
}
