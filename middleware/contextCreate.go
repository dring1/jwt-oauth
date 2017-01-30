package middleware

import (
	"context"
	"net/http"

	"github.com/dring1/jwt-oauth/lib/contextkeys"
)

type ApiContext struct {
	Errors []error
	Value  interface{}
}

func ContextCreate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextkeys.Api, ApiContext{})
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
