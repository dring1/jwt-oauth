package middleware

import (
	"net/http"

	"github.com/dring1/jwt-oauth/token"
)

// TODO: Check if a token is in the blacklist
// if it is return 401
type TokenInvalidatorMiddleware struct {
	TokenService token.Service
	handler      http.Handler
}

func NewTokenInvalidationMiddleware(tokenService token.Service) Middleware {
	return func(next http.Handler) http.Handler {
		return &TokenInvalidatorMiddleware{}
	}
}

func (t *TokenInvalidatorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.TokenService.Validate()

}
