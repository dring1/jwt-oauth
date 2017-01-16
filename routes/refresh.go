package routes

import (
	"log"
	"net/http"

	"github.com/dring1/jwt-oauth/lib/contextkeys"
	"github.com/dring1/jwt-oauth/lib/errors"
	"github.com/dring1/jwt-oauth/token"
)

// given a valid jwt
// generate a new token
// blacklist the token with a TTL until it expires
type RefreshTokenRoute struct {
	Route
	TokenService token.Service `service:"TokenService"`
}

func (rt *RefreshTokenRoute) CompileRoute() (*Route, error) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Check context for token
		iToken := r.Context().Value(contextkeys.Auth)
		tok, ok := (iToken).(token.Token)
		if !ok {
			log.Println("Type assertion failed")
			w.WriteHeader(401)
			errors.ErrorHandler(w, r)
			return
		}
		// Respond with new token with same claims and all
		tokenString, err := rt.TokenService.RefreshToken(&tok)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			errors.ErrorHandler(w, r)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(tokenString))
	}
	rt.Handler = http.HandlerFunc(fn)

	return &rt.Route, nil
}
