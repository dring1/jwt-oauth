package routes

import (
	"github.com/dring1/jwt-oauth/lib/errors"
	"github.com/dring1/jwt-oauth/token"
	"net/http"
)

// given a valid jwt
// generate a new token
// blacklist the token with a TTL until it expires
type RefreshTokenRoute struct {
	Route
	TokenService token.Service
}

const Token = "Token"

func (rt *RefreshTokenRoute) CompileRoute() (*Route, error) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Check context for token
		iToken := r.Context().Value(Token)
		tok, ok := (iToken).(token.Token)
		if !ok {
			w.WriteHeader(500)
			errors.ErrorHandler(w, r)
			return
		}
		// Respond with new token with same claims and all
		tokenString, err := rt.TokenService.RefreshToken(&tok)
		if err != nil {
			// something happened.
		}
		w.WriteHeader(200)
		w.Write([]byte(tokenString))
	}
	rt.Handler = http.HandlerFunc(fn)

	return &rt.Route, nil
}
