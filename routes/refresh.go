package routes

import (
	"net/http"
	"github.com/dring1/jwt-oauth/token"
)

// given a valid jwt
// generate a new token
// blacklist the token with a TTL until it expires
type RefreshTokenRoute struct {
	Route
	TokenService token.Service
}

const Token = "Token"

func (r *RefreshTokenRoute) CompileRoute() (*Route, error) {
	fn := func(w http.ResponseWriter, r *http.Request){
		// have token
		// Here we can assume that it is a valid token ?
		// Check context for token
		tok := r.Context().Value(Token).(token.)
		// Revoke the token
		// Respond with new token with same claims and all
	}
	r.Handler = http.HandlerFunc(fn)
	return &r.Route, nil
} 
