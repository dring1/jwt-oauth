package routes

import (
	"context"
	"net/http"

	jsonresponder "github.com/dring1/jwt-oauth/jsonResponder"
	"github.com/dring1/jwt-oauth/lib/contextkeys"
	"github.com/dring1/jwt-oauth/lib/errors"
	"github.com/dring1/jwt-oauth/token"
)

// given a valid jwt
// generate a new token
// blacklist the token with a TTL until it expires
type RefreshTokenRoute struct {
	Route
	JsonResponder jsonresponder.Service `service:"JsonResponder"`
	TokenService  token.Service         `service:"TokenService"`
}

func (rt *RefreshTokenRoute) CompileRoute() (*Route, error) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Check context for token
		iToken := r.Context().Value(contextkeys.Auth)
		tok, ok := (iToken).(token.Token)
		if !ok {
			w.WriteHeader(401)
			ctx := context.WithValue(r.Context(), contextkeys.Error, errors.InvalidToken)
			r = r.WithContext(ctx)
			rt.JsonResponder.Respond(w, r)
		}
		// Respond with new token with same claims and all
		token, err := rt.TokenService.RefreshToken(&tok)
		if err != nil {
			w.WriteHeader(401)
			ctx := context.WithValue(r.Context(), contextkeys.Error, err)
			r = r.WithContext(ctx)
			rt.JsonResponder.Respond(w, r)
		}

		ctx := context.WithValue(r.Context(), contextkeys.Value, token)
		r = r.WithContext(ctx)
		rt.JsonResponder.Respond(w, r)
	}
	rt.Handler = http.HandlerFunc(fn)

	return &rt.Route, nil
}
