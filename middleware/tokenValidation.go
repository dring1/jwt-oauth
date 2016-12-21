package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/dring1/jwt-oauth/lib/errors"
	"github.com/dring1/jwt-oauth/token"
)

const AuthorizationHeader = "Authorization"
const BearerAuth = "Bearer"

type TokenValidatorMiddleware struct {
	TokenService token.Service
	handler      http.Handler
}

func NewTokenValidationMiddleware(tokenService token.Service) Middleware {

	return func(next http.Handler) http.Handler {
		return &TokenValidatorMiddleware{TokenService: tokenService, handler: next}
	}
}

func (t *TokenValidatorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	header, ok := r.Header[AuthorizationHeader]
	if !ok {
		w.WriteHeader(401)
		w.Write([]byte(errors.InvalidToken))
		return
	}
	tokenHeader := strings.Fields(header[0])
	if len(tokenHeader) != 2 || tokenHeader[0] != BearerAuth {
		w.WriteHeader(401)
		w.Write([]byte(errors.InvalidToken))
		return
	}

	token := tokenHeader[1]
	ok, err := t.TokenService.Validate(token)
	log.Println(err)
	if !ok || err != nil {
		w.WriteHeader(401)
		w.Write([]byte(errors.InvalidToken))
		return
	}
	t.handler.ServeHTTP(w, r)
	return
}
